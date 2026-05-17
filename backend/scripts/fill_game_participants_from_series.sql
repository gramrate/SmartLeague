-- Usage:
-- psql -U user -d db -v game_id="'<uuid>'" -f backend/scripts/fill_game_participants_from_series.sql

BEGIN;

DELETE FROM game_participants
WHERE game_id = :game_id::uuid;

INSERT INTO game_participants (game_id, profile_id)
SELECT :game_id::uuid, x.profile_id
FROM (
  SELECT DISTINCT ON (sp.profile_id)
    sp.profile_id,
    sp.created_at
  FROM series_participants sp
  JOIN games g ON g.series_id = sp.series_id
  WHERE g.id = :game_id::uuid
    AND g.deleted_at IS NULL
  ORDER BY sp.profile_id, sp.created_at ASC
) x
ORDER BY x.created_at ASC;

COMMIT;

