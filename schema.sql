PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  email VARCHAR(70) NOT NULL UNIQUE,
  username VARCHAR(15) NOT NULL UNIQUE,
  first_name VARCHAR(25) NOT NULL,
  last_name VARCHAR(25) NOT NULL,
  dob DATE NOT NULL,
  password VARCHAR(255) NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS posts (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  type VARCHAR(45) NOT NULL,
  title VARCHAR(50) NOT NULL,
  content VARCHAR(500) NOT NULL,
  image_path TEXT NOT NULL,
  user_id INTEGER NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts (user_id);

CREATE TABLE IF NOT EXISTS follows (
  follower_id INTEGER NOT NULL,
  following_id INTEGER NOT NULL,
  PRIMARY KEY (follower_id, following_id),
  FOREIGN KEY (follower_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE,
  FOREIGN KEY (following_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_follows_following_id ON follows (following_id);

CREATE TABLE IF NOT EXISTS profile (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL UNIQUE,
  about VARCHAR(100),
  num_of_followers INTEGER NOT NULL DEFAULT 0,
  num_of_following INTEGER NOT NULL DEFAULT 0,
  num_of_posts INTEGER NOT NULL DEFAULT 0,
  avatar_path TEXT,
  FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS groups (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  title VARCHAR(45) NOT NULL,
  description VARCHAR(300) NOT NULL,
  avatar_path TEXT,
  num_of_users INTEGER NOT NULL DEFAULT 0,
  is_private_chat INTEGER NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS group_users (
  group_id INTEGER NOT NULL,
  user_id INTEGER NOT NULL,
  is_owner INTEGER NOT NULL DEFAULT 0,
  PRIMARY KEY (group_id, user_id),
  FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE,
  FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE ON UPDATE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_group_users_user_id ON group_users (user_id);

CREATE TABLE IF NOT EXISTS events (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  group_id INTEGER NOT NULL,
  creator_id INTEGER NOT NULL,
  title VARCHAR(100) NOT NULL,
  description VARCHAR(500),
  event_datetime DATETIME NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE ON UPDATE CASCADE,
  FOREIGN KEY (creator_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE,
  FOREIGN KEY (group_id, creator_id) REFERENCES group_users (group_id, user_id) ON DELETE CASCADE ON UPDATE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_events_group_id ON events (group_id);
CREATE INDEX IF NOT EXISTS idx_events_creator_id ON events (creator_id);
CREATE INDEX IF NOT EXISTS idx_events_group_creator ON events (group_id, creator_id);

CREATE TABLE IF NOT EXISTS event_options (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  event_id INTEGER NOT NULL,
  option_text VARCHAR(50) NOT NULL,
  FOREIGN KEY (event_id) REFERENCES events (id) ON DELETE CASCADE ON UPDATE CASCADE,
  UNIQUE (event_id, option_text)
);
CREATE INDEX IF NOT EXISTS idx_event_options_event_id ON event_options (event_id);

CREATE TABLE IF NOT EXISTS event_responses (
  event_id INTEGER NOT NULL,
  user_id INTEGER NOT NULL,
  option_id INTEGER NOT NULL,
  PRIMARY KEY (event_id, user_id),
  FOREIGN KEY (event_id) REFERENCES events (id) ON DELETE CASCADE ON UPDATE CASCADE,
  FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE,
  FOREIGN KEY (option_id) REFERENCES event_options (id) ON DELETE CASCADE ON UPDATE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_event_responses_user_id ON event_responses (user_id);
CREATE INDEX IF NOT EXISTS idx_event_responses_option_id ON event_responses (option_id);

CREATE TABLE IF NOT EXISTS chats (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
  group_id INTEGER NOT NULL,
  content VARCHAR(200) NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE ON UPDATE CASCADE,
  FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_chats_group_id ON chats (group_id);
CREATE INDEX IF NOT EXISTS idx_chats_user_id ON chats (user_id);

CREATE TABLE IF NOT EXISTS notifications (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  content VARCHAR(45) NOT NULL,
  type VARCHAR(45) NOT NULL,
  user_id INTEGER NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications (user_id);

CREATE TRIGGER IF NOT EXISTS increase_post_count
AFTER INSERT ON posts

FOR EACH ROW
BEGIN
  UPDATE profile SET num_of_posts = num_of_posts + 1 WHERE user_id = NEW.user_id;
END;

CREATE TRIGGER IF NOT EXISTS decrease_post_count
AFTER DELETE ON posts
FOR EACH ROW
BEGIN
  UPDATE profile SET num_of_posts = num_of_posts - 1 WHERE user_id = OLD.user_id;
END;

CREATE TRIGGER IF NOT EXISTS increase_following_count
AFTER INSERT ON follows
FOR EACH ROW
BEGIN
  UPDATE profile SET num_of_following = num_of_following + 1 WHERE user_id = NEW.follower_id;
END;

CREATE TRIGGER IF NOT EXISTS increase_follower_count
AFTER INSERT ON follows
FOR EACH ROW
BEGIN
  UPDATE profile SET num_of_followers = num_of_followers + 1 WHERE user_id = NEW.following_id;
END;

CREATE TRIGGER IF NOT EXISTS decrease_following_count
AFTER DELETE ON follows
FOR EACH ROW
BEGIN
  UPDATE profile SET num_of_following = num_of_following - 1 WHERE user_id = OLD.follower_id;
END;

CREATE TRIGGER IF NOT EXISTS decrease_follower_count
AFTER DELETE ON follows
FOR EACH ROW
BEGIN
  UPDATE profile SET num_of_followers = num_of_followers - 1 WHERE user_id = OLD.following_id;
END;

CREATE TRIGGER IF NOT EXISTS group_users_AFTER_INSERT
AFTER INSERT ON group_users
FOR EACH ROW
BEGIN
  UPDATE groups SET num_of_users = num_of_users + 1 WHERE id = NEW.group_id;
END;

CREATE TRIGGER IF NOT EXISTS group_users_AFTER_DELETE
AFTER DELETE ON group_users
FOR EACH ROW
BEGIN
  UPDATE groups SET num_of_users = num_of_users - 1 WHERE id = OLD.group_id;
END;