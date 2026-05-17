ALTER TABLE game_results
  ALTER COLUMN yellow_cards TYPE double precision USING yellow_cards::double precision,
  ALTER COLUMN removed TYPE double precision USING removed::double precision;
