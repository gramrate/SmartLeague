-- Guest players (no account required) scoped to a series

CREATE TABLE IF NOT EXISTS series_guests (
	id        uuid PRIMARY KEY,
	series_id uuid NOT NULL,
	nickname  text NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now()
);

DO $$ BEGIN
	ALTER TABLE series_guests
		ADD CONSTRAINT series_guests_series_fk
		FOREIGN KEY (series_id) REFERENCES series(id) ON DELETE CASCADE;
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

CREATE UNIQUE INDEX IF NOT EXISTS series_guests_series_nickname_uq
	ON series_guests(series_id, lower(nickname));

-- Restructure game_results to support guests alongside registered profiles

-- Add surrogate row PK
ALTER TABLE game_results ADD COLUMN IF NOT EXISTS row_id uuid;
UPDATE game_results SET row_id = gen_random_uuid() WHERE row_id IS NULL;
ALTER TABLE game_results ALTER COLUMN row_id SET NOT NULL;
ALTER TABLE game_results ALTER COLUMN row_id SET DEFAULT gen_random_uuid();

-- Add guest_id column
ALTER TABLE game_results ADD COLUMN IF NOT EXISTS guest_id uuid NULL;

DO $$ BEGIN
	ALTER TABLE game_results
		ADD CONSTRAINT game_results_guest_fk
		FOREIGN KEY (guest_id) REFERENCES series_guests(id) ON DELETE CASCADE;
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

-- Swap PK from (game_id, profile_id) to row_id
ALTER TABLE game_results DROP CONSTRAINT IF EXISTS game_results_pkey;
DO $$ BEGIN
	ALTER TABLE game_results ADD PRIMARY KEY (row_id);
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

-- Make profile_id nullable (FK still enforces integrity for non-null values)
ALTER TABLE game_results ALTER COLUMN profile_id DROP NOT NULL;

-- Partial unique indexes replacing the old composite PK
CREATE UNIQUE INDEX IF NOT EXISTS game_results_game_profile_uq
	ON game_results(game_id, profile_id) WHERE profile_id IS NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS game_results_game_guest_uq
	ON game_results(game_id, guest_id) WHERE guest_id IS NOT NULL;

-- Exactly one player identifier must be set
DO $$ BEGIN
	ALTER TABLE game_results ADD CONSTRAINT game_results_player_check CHECK (
		(profile_id IS NOT NULL AND guest_id IS NULL) OR
		(profile_id IS NULL  AND guest_id IS NOT NULL)
	);
EXCEPTION WHEN duplicate_object THEN NULL; END $$;
