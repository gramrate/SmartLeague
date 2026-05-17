ALTER TABLE game_results
	ADD COLUMN IF NOT EXISTS victory_points double precision NOT NULL DEFAULT 0;
