-- Reorder club_state values:
-- old: 0 none, 1 member, 2 leader, 3 president, 4 resident
-- new: 0 none, 1 member, 2 resident, 3 leader, 4 president

UPDATE profiles
SET club_state = 5
WHERE club_state = 4;

UPDATE profiles
SET club_state = 2
WHERE club_state = 5;

UPDATE profiles
SET club_state = 5
WHERE club_state = 2;

UPDATE profiles
SET club_state = 3
WHERE club_state = 5;

UPDATE profiles
SET club_state = 5
WHERE club_state = 3;

UPDATE profiles
SET club_state = 4
WHERE club_state = 5;
