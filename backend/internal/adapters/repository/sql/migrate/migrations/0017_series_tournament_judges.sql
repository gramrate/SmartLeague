ALTER TABLE series ADD COLUMN IF NOT EXISTS is_tournament boolean NOT NULL DEFAULT false;

CREATE TABLE IF NOT EXISTS series_judges (
    series_id  uuid     NOT NULL REFERENCES series(id) ON DELETE CASCADE,
    profile_id uuid     NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    role       smallint NOT NULL DEFAULT 0, -- 0=side, 1=main
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (series_id, profile_id)
);
