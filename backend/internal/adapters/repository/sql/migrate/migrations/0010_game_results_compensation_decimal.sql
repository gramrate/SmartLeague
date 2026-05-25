ALTER TABLE game_results
  ALTER COLUMN compensation TYPE double precision USING compensation::double precision;

