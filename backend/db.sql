CREATE TABLE users (
    username VARCHAR PRIMARY KEY,
    password VARCHAR NOT NULL
);

CREATE TABLE rooms (
    room_id VARCHAR PRIMARY KEY,
    host VARCHAR,
    FOREIGN KEY (host) REFERENCES users(username)
);

CREATE TABLE messages (
    message_id VARCHAR PRIMARY KEY,
    room_id VARCHAR,
    timestamp TIMESTAMPTZ NOT NULL,
    author VARCHAR,
    content TEXT,
    FOREIGN KEY (room_id) REFERENCES rooms(room_id),
    FOREIGN KEY (author) REFERENCES users(username) ON DELETE SET NULL;
);

CREATE INDEX idx_messages_room_id_timestamp ON messages(room_id, timestamp);

CREATE TABLE users_rooms (
    username VARCHAR,
    room_id VARCHAR,
    PRIMARY KEY (username, room_id),
    FOREIGN KEY (username) REFERENCES users(username),
    FOREIGN KEY (room_id) REFERENCES rooms(room_id)
);

CREATE INDEX idx_users_rooms_room_id ON users_rooms(room_id);
