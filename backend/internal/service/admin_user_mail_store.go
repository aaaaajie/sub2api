package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/google/uuid"
)

func (s *AdminUserMailService) insertJob(ctx context.Context, actorID int64, id uuid.UUID, subject, html string, isTest bool, recipients []AdminUserMailRecipient) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	result, err := tx.ExecContext(ctx, `
		INSERT INTO admin_user_mail_jobs (id, created_by, subject, html_body, is_test) VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO NOTHING`, id, actorID, subject, html, isTest)
	if err != nil {
		return err
	}
	inserted, err := result.RowsAffected()
	if err != nil || inserted == 0 {
		return err
	}
	for _, recipient := range recipients {
		if _, err := tx.ExecContext(ctx, `INSERT INTO admin_user_mail_recipients (job_id, user_id, email) VALUES ($1, $2, $3)`,
			id, recipient.UserID, recipient.Email); err != nil {
			return err
		}
	}
	return tx.Commit()
}

const adminUserMailSummarySelect = `
	SELECT j.id, j.subject, j.is_test, j.created_at,
		count(*) FILTER (WHERE r.status IN ('pending','sending')),
		count(*) FILTER (WHERE r.status = 'accepted'), count(*) FILTER (WHERE r.status = 'failed'),
		count(*) FILTER (WHERE r.status = 'skipped'), count(*) FILTER (WHERE r.status = 'uncertain')
	FROM admin_user_mail_jobs j JOIN admin_user_mail_recipients r ON r.job_id = j.id`

func (s *AdminUserMailService) List(ctx context.Context, page pagination.PaginationParams) ([]AdminUserMailJob, int64, error) {
	var total int64
	if err := s.db.QueryRowContext(ctx, `SELECT count(*) FROM admin_user_mail_jobs`).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := s.db.QueryContext(ctx, adminUserMailSummarySelect+`
		GROUP BY j.id ORDER BY j.created_at DESC, j.id LIMIT $1 OFFSET $2`, page.Limit(), page.Offset())
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	jobs := make([]AdminUserMailJob, 0)
	for rows.Next() {
		var job AdminUserMailJob
		if err := rows.Scan(&job.ID, &job.Subject, &job.IsTest, &job.CreatedAt,
			&job.Pending, &job.Accepted, &job.Failed, &job.Skipped, &job.Uncertain); err != nil {
			return nil, 0, err
		}
		jobs = append(jobs, job)
	}
	return jobs, total, rows.Err()
}

func (s *AdminUserMailService) Get(ctx context.Context, id uuid.UUID) (AdminUserMailJobDetail, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT j.id, j.subject, j.is_test, j.created_at, r.user_id, r.email, r.status, r.reason
		FROM admin_user_mail_jobs j JOIN admin_user_mail_recipients r ON r.job_id = j.id
		WHERE j.id = $1 ORDER BY r.user_id`, id)
	if err != nil {
		return AdminUserMailJobDetail{}, err
	}
	defer rows.Close()
	var job AdminUserMailJobDetail
	for rows.Next() {
		var recipient AdminUserMailRecipient
		if err := rows.Scan(&job.ID, &job.Subject, &job.IsTest, &job.CreatedAt,
			&recipient.UserID, &recipient.Email, &recipient.Status, &recipient.Reason); err != nil {
			return AdminUserMailJobDetail{}, err
		}
		job.Recipients = append(job.Recipients, recipient)
		switch recipient.Status {
		case "pending", "sending":
			job.Pending++
		case "accepted":
			job.Accepted++
		case "failed":
			job.Failed++
		case "skipped":
			job.Skipped++
		case "uncertain":
			job.Uncertain++
		}
	}
	if err := rows.Err(); err != nil {
		return AdminUserMailJobDetail{}, err
	}
	if len(job.Recipients) == 0 {
		return AdminUserMailJobDetail{}, errUserMailNotFound
	}
	return job, nil
}

func (s *AdminUserMailService) Retry(ctx context.Context, id uuid.UUID, userID int64) (AdminUserMailJobDetail, error) {
	_, err := s.db.ExecContext(ctx, `UPDATE admin_user_mail_recipients
		SET status = 'pending', started_at = NULL, reason = '' WHERE job_id = $1 AND user_id = $2 AND status = 'failed'`, id, userID)
	if err != nil {
		return AdminUserMailJobDetail{}, err
	}
	return s.Get(ctx, id)
}

type adminUserMailAttempt struct {
	JobID   uuid.UUID
	UserID  int64
	Email   string
	Subject string
	HTML    string
	Started time.Time
}

func (s *AdminUserMailService) claim(ctx context.Context) (*adminUserMailAttempt, error) {
	var attempt adminUserMailAttempt
	err := s.db.QueryRowContext(ctx, `
		WITH expired AS (
			UPDATE admin_user_mail_recipients SET status = 'uncertain', reason = 'interrupted'
			WHERE status = 'sending' AND started_at < now() - interval '10 minutes'
		), next AS (
			SELECT r.job_id, r.user_id, j.subject, j.html_body
			FROM admin_user_mail_recipients r JOIN admin_user_mail_jobs j ON j.id = r.job_id
			WHERE r.status = 'pending' ORDER BY j.created_at, j.id, r.user_id LIMIT 1 FOR UPDATE OF r SKIP LOCKED
		)
		UPDATE admin_user_mail_recipients r SET status = 'sending', started_at = now()
		FROM next WHERE r.job_id = next.job_id AND r.user_id = next.user_id
		RETURNING r.job_id, r.user_id, r.email, next.subject, next.html_body, r.started_at`).
		Scan(&attempt.JobID, &attempt.UserID, &attempt.Email, &attempt.Subject, &attempt.HTML, &attempt.Started)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &attempt, err
}
