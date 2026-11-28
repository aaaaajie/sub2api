package service

import (
	"context"
	"database/sql"
	"errors"
	"net/mail"
	"strings"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

const maxUserMailRecipients = 500

var (
	errUserMailAudienceLimit   = infraerrors.BadRequest("USER_MAIL_AUDIENCE_LIMIT", "Select at most 500 users per send.")
	errUserMailEmptyAudience   = infraerrors.BadRequest("USER_MAIL_EMPTY_AUDIENCE", "No valid recipient email addresses were found.")
	errUserMailAudienceChanged = infraerrors.Conflict("USER_MAIL_AUDIENCE_CHANGED", "The recipient list changed. Preview the message again before sending.")
	errUserMailNotFound        = infraerrors.NotFound("USER_MAIL_NOT_FOUND", "Email job not found.")
)

type AdminUserMailService struct {
	db       *sql.DB
	users    UserRepository
	email    *EmailService
	settings *SettingService
	cancel   context.CancelFunc
	done     chan struct{}
}

func ProvideAdminUserMailService(db *sql.DB, users UserRepository, email *EmailService, settings *SettingService) *AdminUserMailService {
	ctx, cancel := context.WithCancel(context.Background())
	s := &AdminUserMailService{db: db, users: users, email: email, settings: settings, cancel: cancel, done: make(chan struct{})}
	go s.run(ctx)
	return s
}

func (s *AdminUserMailService) Stop(ctx context.Context) error {
	s.cancel()
	select {
	case <-s.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

type AdminUserMailPreview struct {
	HTML    string  `json:"html"`
	UserIDs []int64 `json:"user_ids"`
	Skipped int     `json:"skipped"`
}

type AdminUserMailRecipient struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	Status string `json:"status"`
	Reason string `json:"reason,omitempty"`
}

type AdminUserMailJob struct {
	ID        uuid.UUID `json:"id"`
	Subject   string    `json:"subject"`
	IsTest    bool      `json:"is_test"`
	CreatedAt time.Time `json:"created_at"`
	Pending   int       `json:"pending"`
	Accepted  int       `json:"accepted"`
	Failed    int       `json:"failed"`
	Skipped   int       `json:"skipped"`
	Uncertain int       `json:"uncertain"`
}

type AdminUserMailJobDetail struct {
	AdminUserMailJob
	Recipients []AdminUserMailRecipient `json:"recipients"`
}

func (s *AdminUserMailService) Preview(ctx context.Context, draft AdminUserMailDraft, ids []int64, filters *UserListFilters) (AdminUserMailPreview, error) {
	_, html, err := s.render(ctx, draft)
	if err != nil {
		return AdminUserMailPreview{}, err
	}
	recipients, skipped, err := s.recipients(ctx, ids, filters)
	if err != nil {
		return AdminUserMailPreview{}, err
	}
	preview := AdminUserMailPreview{HTML: html, UserIDs: make([]int64, 0, len(recipients)), Skipped: skipped}
	for _, recipient := range recipients {
		preview.UserIDs = append(preview.UserIDs, recipient.UserID)
	}
	return preview, nil
}

func (s *AdminUserMailService) recipients(ctx context.Context, ids []int64, filters *UserListFilters) ([]AdminUserMailRecipient, int, error) {
	if len(ids) > maxUserMailRecipients {
		return nil, 0, errUserMailAudienceLimit
	}
	candidates := make([]AdminUserMailRecipient, 0)
	requested := len(ids)
	if filters != nil {
		noSubscriptions := false
		listFilters := *filters
		listFilters.IncludeSubscriptions = &noSubscriptions
		users, page, err := s.users.ListWithFilters(ctx, pagination.PaginationParams{Page: 1, PageSize: maxUserMailRecipients, SortBy: "id", SortOrder: "asc"}, listFilters)
		if err != nil {
			return nil, 0, err
		}
		if page.Total > maxUserMailRecipients {
			return nil, 0, errUserMailAudienceLimit
		}
		requested = len(users)
		for _, user := range users {
			candidates = append(candidates, AdminUserMailRecipient{UserID: user.ID, Email: user.Email})
		}
	} else if len(ids) > 0 {
		rows, err := s.db.QueryContext(ctx, `SELECT id, email FROM users WHERE deleted_at IS NULL AND id = ANY($1) ORDER BY id`, pq.Array(ids))
		if err != nil {
			return nil, 0, err
		}
		defer rows.Close()
		for rows.Next() {
			var recipient AdminUserMailRecipient
			if err := rows.Scan(&recipient.UserID, &recipient.Email); err != nil {
				return nil, 0, err
			}
			candidates = append(candidates, recipient)
		}
		if err := rows.Err(); err != nil {
			return nil, 0, err
		}
	}
	recipients := eligibleUserMailRecipients(candidates)
	return recipients, requested - len(recipients), nil
}

func eligibleUserMailRecipients(candidates []AdminUserMailRecipient) []AdminUserMailRecipient {
	seen := make(map[string]struct{}, len(candidates))
	recipients := make([]AdminUserMailRecipient, 0, len(candidates))
	for _, recipient := range candidates {
		recipient.Email = strings.TrimSpace(recipient.Email)
		normalized := strings.ToLower(recipient.Email)
		address, err := mail.ParseAddress(recipient.Email)
		_, duplicate := seen[normalized]
		if err != nil || address.Address != recipient.Email || strings.HasSuffix(normalized, ".invalid") || duplicate {
			continue
		}
		seen[normalized] = struct{}{}
		recipients = append(recipients, recipient)
	}
	return recipients
}

func (s *AdminUserMailService) Queue(ctx context.Context, actorID int64, id uuid.UUID, draft AdminUserMailDraft, ids []int64, isTest bool) (AdminUserMailJob, error) {
	existing, err := s.Get(ctx, id)
	if err == nil {
		return existing.AdminUserMailJob, nil
	}
	if !errors.Is(err, errUserMailNotFound) {
		return AdminUserMailJob{}, err
	}
	subject, html, err := s.render(ctx, draft)
	if err != nil {
		return AdminUserMailJob{}, err
	}
	if isTest {
		ids = []int64{actorID}
	}
	recipients, _, err := s.recipients(ctx, ids, nil)
	if err != nil {
		return AdminUserMailJob{}, err
	}
	if len(recipients) == 0 {
		return AdminUserMailJob{}, errUserMailEmptyAudience
	}
	if len(recipients) != len(ids) {
		return AdminUserMailJob{}, errUserMailAudienceChanged
	}
	config, err := s.email.GetSMTPConfig(ctx)
	if err != nil {
		return AdminUserMailJob{}, err
	}
	if config.From == "" {
		return AdminUserMailJob{}, infraerrors.BadRequest("USER_MAIL_SMTP_SENDER_REQUIRED", "Configure an SMTP sender address before sending.")
	}
	if err := s.insertJob(ctx, actorID, id, subject, html, isTest, recipients); err != nil {
		return AdminUserMailJob{}, err
	}
	job, err := s.Get(ctx, id)
	return job.AdminUserMailJob, err
}
