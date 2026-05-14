ALTER TABLE game_results
	ALTER COLUMN best_move DROP DEFAULT,
	ALTER COLUMN removed DROP DEFAULT;

ALTER TABLE game_results
	ALTER COLUMN best_move TYPE text USING CASE WHEN best_move THEN '1' ELSE NULL END,
	ALTER COLUMN compensation TYPE double precision USING compensation::double precision,
	ALTER COLUMN removed TYPE integer USING CASE WHEN removed THEN 1 ELSE 0 END,
	ALTER COLUMN extra_points TYPE double precision USING extra_points::double precision,
	ALTER COLUMN total_points TYPE double precision USING total_points::double precision;

ALTER TABLE game_results
	ALTER COLUMN best_move SET DEFAULT NULL,
	ALTER COLUMN removed SET DEFAULT 0;
