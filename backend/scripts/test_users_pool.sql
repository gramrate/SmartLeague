-- Test pool: 10 users in club + join series
-- Club ID / Series ID from task:
-- 22ca91b8-644b-4e33-af9d-3d38de20668f
--
-- Password hash below is a stub. If you need login for these users,
-- replace `password_hash` with a valid bcrypt hash.

WITH pool(id, nickname, name, email) AS (
	VALUES
		('7f5cf03e-6f2a-41c8-b420-c86a79db8771'::uuid, 'test_user_01', 'Test User 01', 'test_user_01@smartleague.local'),
		('6726cfc2-dd1e-4b99-b684-519edf08ec29'::uuid, 'test_user_02', 'Test User 02', 'test_user_02@smartleague.local'),
		('7f3599d8-89ca-49f6-9ddf-4693d11de8e6'::uuid, 'test_user_03', 'Test User 03', 'test_user_03@smartleague.local'),
		('f90a5964-4f4e-4fed-b146-7c533674f6c1'::uuid, 'test_user_04', 'Test User 04', 'test_user_04@smartleague.local'),
		('22df3a92-1ca0-494d-b2bd-30bd2164535d'::uuid, 'test_user_05', 'Test User 05', 'test_user_05@smartleague.local'),
		('b8f44795-1804-4264-a861-31d91f23f39f'::uuid, 'test_user_06', 'Test User 06', 'test_user_06@smartleague.local'),
		('6ec9b4b7-a7b9-4600-8dae-c5f7f55fdd4a'::uuid, 'test_user_07', 'Test User 07', 'test_user_07@smartleague.local'),
		('5f59a19a-e708-42ca-92ef-389e2f051737'::uuid, 'test_user_08', 'Test User 08', 'test_user_08@smartleague.local'),
		('7f0b2476-8f4b-4f98-b44c-f9f1bd7ac426'::uuid, 'test_user_09', 'Test User 09', 'test_user_09@smartleague.local'),
		('59198d7d-c8e0-4cd6-952c-8acddbe2ac27'::uuid, 'test_user_10', 'Test User 10', 'test_user_10@smartleague.local')
)
INSERT INTO profiles (id, nickname, name, show_name, description, email, password_hash, club_id, club_state, role)
SELECT
	p.id,
	p.nickname,
	p.name,
	true,
	NULL,
	p.email,
	'test-password-hash',
	'22ca91b8-644b-4e33-af9d-3d38de20668f'::uuid,
	1,
	0
FROM pool p
ON CONFLICT (id) DO UPDATE
SET
	nickname = EXCLUDED.nickname,
	name = EXCLUDED.name,
	show_name = EXCLUDED.show_name,
	description = EXCLUDED.description,
	email = EXCLUDED.email,
	club_id = EXCLUDED.club_id,
	club_state = EXCLUDED.club_state,
	role = EXCLUDED.role,
	updated_at = now();

WITH ids AS (
	SELECT id
	FROM profiles
	WHERE email LIKE 'test_user_%@smartleague.local'
)
INSERT INTO series_participants (series_id, profile_id)
SELECT '22ca91b8-644b-4e33-af9d-3d38de20668f'::uuid, ids.id
FROM ids
WHERE EXISTS (
	SELECT 1 FROM series WHERE id = '22ca91b8-644b-4e33-af9d-3d38de20668f'::uuid
)
ON CONFLICT DO NOTHING;
