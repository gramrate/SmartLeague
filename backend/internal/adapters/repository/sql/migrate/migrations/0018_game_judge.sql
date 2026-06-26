ALTER TABLE games ADD COLUMN IF NOT EXISTS game_judge_id uuid REFERENCES profiles(id) ON DELETE SET NULL;
ALTER TABLE games ADD COLUMN IF NOT EXISTS game_judge_confirmed boolean NOT NULL DEFAULT false;
