-- Auth-related additions

ALTER TABLE profiles ADD COLUMN IF NOT EXISTS surname text NOT NULL DEFAULT '';

CREATE TABLE IF NOT EXISTS refresh_tokens (
	id uuid PRIMARY KEY,
	user_id uuid NOT NULL UNIQUE,
	jti text NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now()
);

DO $$
BEGIN
	ALTER TABLE refresh_tokens
		ADD CONSTRAINT refresh_tokens_user_fk
		FOREIGN KEY (user_id) REFERENCES profiles(id) ON DELETE CASCADE;
EXCEPTION
	WHEN duplicate_object THEN NULL;
END $$;

