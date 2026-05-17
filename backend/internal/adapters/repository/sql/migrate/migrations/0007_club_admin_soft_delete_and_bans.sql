-- Club bans table
CREATE TABLE IF NOT EXISTS club_bans (
  club_id uuid NOT NULL,
  profile_id uuid NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (club_id, profile_id)
);

DO $$
BEGIN
  ALTER TABLE club_bans
    ADD CONSTRAINT club_bans_club_fk
    FOREIGN KEY (club_id) REFERENCES clubs(id) ON DELETE CASCADE;
EXCEPTION
  WHEN duplicate_object THEN NULL;
END $$;

DO $$
BEGIN
  ALTER TABLE club_bans
    ADD CONSTRAINT club_bans_profile_fk
    FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE;
EXCEPTION
  WHEN duplicate_object THEN NULL;
END $$;

-- Soft-delete markers
ALTER TABLE series ADD COLUMN IF NOT EXISTS deleted_at timestamptz NULL;
ALTER TABLE games ADD COLUMN IF NOT EXISTS deleted_at timestamptz NULL;

CREATE INDEX IF NOT EXISTS series_deleted_at_idx ON series (deleted_at);
CREATE INDEX IF NOT EXISTS games_deleted_at_idx ON games (deleted_at);
