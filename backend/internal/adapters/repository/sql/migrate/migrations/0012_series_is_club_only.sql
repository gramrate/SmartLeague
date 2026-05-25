ALTER TABLE series
	ADD COLUMN IF NOT EXISTS is_club_only boolean NOT NULL DEFAULT false;
