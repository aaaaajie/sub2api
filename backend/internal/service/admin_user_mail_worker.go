package service

import (
	"context"
	"errors"
	"log/slog"
	"net/textproto"
	"time"
)

func (s *AdminUserMailService) run(ctx context.Context) {
	defer close(s.done)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for ctx.Err() == nil {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			lookup, cancel := context.WithTimeout(ctx, 10*time.Second)
			attempt, err := s.claim(lookup)
			cancel()
			if err != nil {
				if ctx.Err() == nil {
					slog.Warn("admin user email claim failed", "error", err)
				}
				continue
			}
			if attempt != nil {
				s.deliver(ctx, *attempt)
				ticker.Reset(2 * time.Second)
			}
		}
	}
}

func (s *AdminUserMailService) deliver(ctx context.Context, attempt adminUserMailAttempt) {
	lookup, cancel := context.WithTimeout(ctx, 10*time.Second)
	var current bool
	err := s.db.QueryRowContext(lookup, `SELECT EXISTS(
		SELECT 1 FROM users WHERE id = $1 AND deleted_at IS NULL AND btrim(email) = $2
	)`, attempt.UserID, attempt.Email).Scan(&current)
	cancel()
	status, reason := "accepted", ""
	switch {
	case err != nil:
		status, reason = "failed", "recipient_lookup"
	case !current:
		status, reason = "skipped", "recipient_changed"
	default:
		status, reason = adminUserMailDeliveryResult(s.email.SendEmail(ctx, attempt.Email, attempt.Subject, attempt.HTML))
	}
	// Record the SMTP outcome even if graceful shutdown has begun.
	update, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
	defer cancel()
	if _, err := s.db.ExecContext(update, `UPDATE admin_user_mail_recipients SET status = $4, reason = $5
		WHERE job_id = $1 AND user_id = $2 AND started_at = $3 AND status = 'sending'`,
		attempt.JobID, attempt.UserID, attempt.Started, status, reason); err != nil {
		slog.Error("admin user email outcome persistence failed", "job_id", attempt.JobID, "user_id", attempt.UserID)
	}
}

func adminUserMailDeliveryResult(err error) (string, string) {
	if err == nil {
		return "accepted", ""
	}
	if errors.Is(err, errSMTPAcceptanceUnknown) {
		return "uncertain", "smtp_confirmation"
	}
	var smtpError *textproto.Error
	if errors.As(err, &smtpError) {
		return "failed", "smtp_rejected"
	}
	return "failed", "smtp_connection"
}
