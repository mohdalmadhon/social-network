-- group id refrence to if the post is 
--  public: group id == null
CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    content VARCHAR(1000) NOT NULL,
    image_path TEXT,
    allow_comments INTEGER NOT NULL DEFAULT 1,
    location TEXT,
    group_id INTEGER,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE,
    FOREIGN KEY (group_id) REFERENCES user_posts_groups(id) ON DELETE
    SET
        NULL
);

CREATE TABLE IF NOT EXISTS user_posts_groups (
    id INTEGER PRIMARY KEY,
    name VARCHAR(15) NOT NULL,
    user_id INTEGER NOT NULL,
    users TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE,
    UNIQUE (user_id, name),
    UNIQUE (user_id, users)
);

CREATE TABLE IF NOT EXISTS post_user_tags (
    user_id INTEGER NOT NULL,
    post_id INTEGER NOT NULL,
    PRIMARY KEY (user_id, post_id),
    FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
);

CREATE TRIGGER IF NOT EXISTS trg_increase_post
AFTER
INSERT
    ON posts FOR EACH ROW BEGIN
UPDATE
    profile
SET
    num_of_posts = num_of_posts + 1
WHERE
    user_id = NEW.user_id;

END;

CREATE TRIGGER IF NOT EXISTS trg_decrease_post
AFTER
    DELETE ON posts FOR EACH ROW BEGIN
UPDATE
    profile
SET
    num_of_posts = num_of_posts - 1
WHERE
    user_id = OLD.user_id;

END;

CREATE TRIGGER IF NOT EXISTS trg_post_groups
AFTER
INSERT
    ON user FOR EACH ROW BEGIN
INSERT INTO
    user_posts_groups (id, name, user_id, users)
VALUES
    (0, 'public', NULL, '');

INSERT INTO
    user_posts_groups (id, name, user_id, users)
VALUES
    (-1, 'private', NULL, '');

END;

