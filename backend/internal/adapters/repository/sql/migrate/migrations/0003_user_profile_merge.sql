-- Merge profile fields into user model, remove deprecated surname

ALTER TABLE profiles DROP COLUMN IF EXISTS surname;
