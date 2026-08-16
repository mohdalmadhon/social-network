
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

CREATE TABLE IF NOT EXISTS groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title VARCHAR(45) NOT NULL,
    description VARCHAR(500) NOT NULL
);

CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    type VARCHAR(45) NOT NULL,
    title VARCHAR(50) NOT NULL,
    content VARCHAR(500) NOT NULL,
    image_path TEXT NOT NULL,
    user_id INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    group_id INTEGER NOT NULL,

    like_count INTEGER NOT NULL DEFAULT 0,
    dislike_count INTEGER NOT NULL DEFAULT 0,
    comment_count INTEGER NOT NULL DEFAULT 0,

    FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    FOREIGN KEY (group_id)
        REFERENCES groups(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS follows (
    follower_id INTEGER NOT NULL,
    following_id INTEGER NOT NULL,

    PRIMARY KEY (follower_id, following_id),

    FOREIGN KEY (follower_id)
        REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    FOREIGN KEY (following_id)
        REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);


CREATE TABLE IF NOT EXISTS profile (
    user_id INTEGER PRIMARY KEY,

    about VARCHAR(100),
    num_of_followers INTEGER NOT NULL DEFAULT 0,
    num_of_following INTEGER NOT NULL DEFAULT 0,
    num_of_posts INTEGER NOT NULL DEFAULT 0,
    avatar_path TEXT,

    FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS chats (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    type VARCHAR(45) NOT NULL,
    group_id INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    num_of_members INTEGER NOT NULL DEFAULT 0,

    FOREIGN KEY (group_id)
        REFERENCES groups(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    sender_id INTEGER NOT NULL,
    chat_id INTEGER NOT NULL,

    content TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (sender_id)
        REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    FOREIGN KEY (chat_id)
        REFERENCES chats(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);


CREATE TABLE IF NOT EXISTS chat_users (
    user_id INTEGER NOT NULL,
    chat_id INTEGER NOT NULL,
    is_owner INTEGER NOT NULL DEFAULT 0,

    PRIMARY KEY (user_id, chat_id),

    FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    FOREIGN KEY (chat_id)
        REFERENCES chats(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);


CREATE TABLE IF NOT EXISTS comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    user_id INTEGER NOT NULL,
    post_id INTEGER NOT NULL,

    content VARCHAR(200) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    like_count INTEGER NOT NULL DEFAULT 0,
    dislike_count INTEGER NOT NULL DEFAULT 0,

    FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    FOREIGN KEY (post_id)
        REFERENCES posts(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    title VARCHAR(50) NOT NULL,
    content VARCHAR(500) NOT NULL,

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    group_id INTEGER NOT NULL,
    creator_id INTEGER NOT NULL,

    FOREIGN KEY (group_id)
        REFERENCES groups(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    FOREIGN KEY (creator_id)
        REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS votes_options (
    id INTEGER PRIMARY KEY,
    option VARCHAR(45) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS event_votes (
    event_id INTEGER NOT NULL,
    vote_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,

    PRIMARY KEY (event_id, user_id),

    FOREIGN KEY (event_id)
        REFERENCES events(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    FOREIGN KEY (vote_id)
        REFERENCES votes_options(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS post_reactions (
    user_id INTEGER NOT NULL,
    post_id INTEGER NOT NULL,
    value INTEGER NOT NULL,

    PRIMARY KEY (user_id, post_id),

    FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    FOREIGN KEY (post_id)
        REFERENCES posts(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS comments_reactions (
    comment_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    value INTEGER NOT NULL,

    PRIMARY KEY (comment_id, user_id),

    FOREIGN KEY (comment_id)
        REFERENCES comments(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

