-- Core domain tables

CREATE TABLE IF NOT EXISTS profiles (
	id uuid PRIMARY KEY,
	nickname text NOT NULL DEFAULT '',
	name text NOT NULL,
	show_name boolean NOT NULL DEFAULT true,
	description text NULL,
	email text NOT NULL,
	password_hash text NOT NULL,
	club_id uuid NULL,
	club_state smallint NOT NULL DEFAULT 0,
	role smallint NOT NULL DEFAULT 0,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX IF NOT EXISTS profiles_email_uq ON profiles (lower(email));

CREATE TABLE IF NOT EXISTS clubs (
	id uuid PRIMARY KEY,
	creator_id uuid NOT NULL,
	name text NOT NULL,
	description text NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now()
);

DO $$
BEGIN
	ALTER TABLE clubs
		ADD CONSTRAINT clubs_creator_id_fk
		FOREIGN KEY (creator_id) REFERENCES profiles(id) ON DELETE RESTRICT;
EXCEPTION
	WHEN duplicate_object THEN NULL;
END $$;

DO $$
BEGIN
	ALTER TABLE profiles
		ADD CONSTRAINT profiles_club_id_fk
		FOREIGN KEY (club_id) REFERENCES clubs(id) ON DELETE SET NULL;
EXCEPTION
	WHEN duplicate_object THEN NULL;
END $$;

-- Series / Games

CREATE TABLE IF NOT EXISTS series (
	id uuid PRIMARY KEY,
	club_id uuid NOT NULL,
	creator_id uuid NOT NULL,
	name text NOT NULL,
	scoring_rules text NOT NULL,
	start_at timestamptz NOT NULL,
	end_at timestamptz NOT NULL,
	description text NULL,
	price_rub integer NOT NULL DEFAULT 0,
	is_closed boolean NOT NULL DEFAULT false,
	game_type smallint NOT NULL DEFAULT 0,
	status smallint NOT NULL DEFAULT 0,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now()
);

DO $$
BEGIN
	ALTER TABLE series
		ADD CONSTRAINT series_club_id_fk
		FOREIGN KEY (club_id) REFERENCES clubs(id) ON DELETE CASCADE;
EXCEPTION
	WHEN duplicate_object THEN NULL;
END $$;

DO $$
BEGIN
	ALTER TABLE series
		ADD CONSTRAINT series_creator_id_fk
		FOREIGN KEY (creator_id) REFERENCES profiles(id) ON DELETE RESTRICT;
EXCEPTION
	WHEN duplicate_object THEN NULL;
END $$;

CREATE INDEX IF NOT EXISTS series_club_id_idx ON series (club_id);

CREATE TABLE IF NOT EXISTS series_participants (
	series_id uuid NOT NULL,
	profile_id uuid NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY (series_id, profile_id)
);

DO $$
BEGIN
	ALTER TABLE series_participants
		ADD CONSTRAINT series_participants_series_fk
		FOREIGN KEY (series_id) REFERENCES series(id) ON DELETE CASCADE;
EXCEPTION
	WHEN duplicate_object THEN NULL;
END $$;

DO $$
BEGIN
	ALTER TABLE series_participants
		ADD CONSTRAINT series_participants_profile_fk
		FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE;
EXCEPTION
	WHEN duplicate_object THEN NULL;
END $$;

CREATE TABLE IF NOT EXISTS games (
	id uuid PRIMARY KEY,
	series_id uuid NOT NULL,
	name text NOT NULL,
	number integer NOT NULL,
	description text NULL,
	host_id uuid NULL,
	status smallint NOT NULL DEFAULT 0,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	UNIQUE (series_id, number)
);

DO $$
BEGIN
	ALTER TABLE games
		ADD CONSTRAINT games_series_fk
		FOREIGN KEY (series_id) REFERENCES series(id) ON DELETE CASCADE;
EXCEPTION
	WHEN duplicate_object THEN NULL;
END $$;

DO $$
BEGIN
	ALTER TABLE games
		ADD CONSTRAINT games_host_fk
		FOREIGN KEY (host_id) REFERENCES profiles(id) ON DELETE SET NULL;
EXCEPTION
	WHEN duplicate_object THEN NULL;
END $$;

CREATE INDEX IF NOT EXISTS games_series_idx ON games (series_id);

CREATE TABLE IF NOT EXISTS game_participants (
	game_id uuid NOT NULL,
	profile_id uuid NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY (game_id, profile_id)
);

DO $$
BEGIN
	ALTER TABLE game_participants
		ADD CONSTRAINT game_participants_game_fk
		FOREIGN KEY (game_id) REFERENCES games(id) ON DELETE CASCADE;
EXCEPTION
	WHEN duplicate_object THEN NULL;
END $$;

DO $$
BEGIN
	ALTER TABLE game_participants
		ADD CONSTRAINT game_participants_profile_fk
		FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE;
EXCEPTION
	WHEN duplicate_object THEN NULL;
END $$;

CREATE TABLE IF NOT EXISTS game_results (
	game_id uuid NOT NULL,
	profile_id uuid NOT NULL,
	place integer NULL,
	role text NULL,
	best_move text NULL,
	first_killed boolean NOT NULL DEFAULT false,
	compensation double precision NOT NULL DEFAULT 0,
	yellow_cards integer NOT NULL DEFAULT 0,
	removed integer NOT NULL DEFAULT 0,
	extra_points double precision NOT NULL DEFAULT 0,
	total_points double precision NOT NULL DEFAULT 0,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY (game_id, profile_id)
);

DO $$
BEGIN
	ALTER TABLE game_results
		ADD CONSTRAINT game_results_game_fk
		FOREIGN KEY (game_id) REFERENCES games(id) ON DELETE CASCADE;
EXCEPTION
	WHEN duplicate_object THEN NULL;
END $$;

DO $$
BEGIN
	ALTER TABLE game_results
		ADD CONSTRAINT game_results_profile_fk
		FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE;
EXCEPTION
	WHEN duplicate_object THEN NULL;
END $$;
