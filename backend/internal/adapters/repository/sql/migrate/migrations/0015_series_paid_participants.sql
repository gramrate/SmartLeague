CREATE TABLE IF NOT EXISTS series_paid_participants (
    series_id uuid NOT NULL,
    profile_id uuid NOT NULL,
    paid_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (series_id, profile_id)
);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'series_paid_participants_series_fk'
    ) THEN
        ALTER TABLE series_paid_participants
            ADD CONSTRAINT series_paid_participants_series_fk
            FOREIGN KEY (series_id) REFERENCES series(id) ON DELETE CASCADE;
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'series_paid_participants_profile_fk'
    ) THEN
        ALTER TABLE series_paid_participants
            ADD CONSTRAINT series_paid_participants_profile_fk
            FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE;
    END IF;
END $$;
