CREATE TABLE IF NOT EXISTS game_drafts (
    game_id         uuid PRIMARY KEY REFERENCES games(id) ON DELETE CASCADE,
    rows            jsonb NOT NULL DEFAULT '[]',
    judge_id        uuid REFERENCES profiles(id) ON DELETE SET NULL,
    judge_confirmed boolean NOT NULL DEFAULT false,
    updated_at      timestamptz NOT NULL DEFAULT now()
);
