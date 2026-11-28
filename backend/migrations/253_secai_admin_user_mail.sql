CREATE TABLE IF NOT EXISTS admin_user_mail_jobs (
    id uuid PRIMARY KEY,
    created_by bigint NOT NULL,
    subject text NOT NULL,
    html_body text NOT NULL,
    is_test boolean NOT NULL DEFAULT false,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS admin_user_mail_jobs_created_at_idx
    ON admin_user_mail_jobs (created_at DESC, id);

CREATE TABLE IF NOT EXISTS admin_user_mail_recipients (
    job_id uuid NOT NULL REFERENCES admin_user_mail_jobs(id) ON DELETE CASCADE,
    user_id bigint NOT NULL,
    email text NOT NULL,
    status text NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'sending', 'accepted', 'failed', 'skipped', 'uncertain')),
    reason text NOT NULL DEFAULT '',
    started_at timestamptz,
    PRIMARY KEY (job_id, user_id),
    UNIQUE (job_id, email)
);

CREATE INDEX IF NOT EXISTS admin_user_mail_recipients_pending_idx
    ON admin_user_mail_recipients (job_id, user_id) WHERE status = 'pending';

CREATE INDEX IF NOT EXISTS admin_user_mail_recipients_sending_idx
    ON admin_user_mail_recipients (started_at) WHERE status = 'sending';
