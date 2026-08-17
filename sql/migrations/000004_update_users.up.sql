/*
    updated users table to add "updated_at" column and it main job is to stop spamming
    the database with changes which can slow down the server. I assume that I will or other person
    will add a check on the last updated time before actully allowing the user to update
*/

ALTER TABLE users
ADD COLUMN updated_at DATETIME 
DEFAULT CURRENT_TIMESTAMP;
