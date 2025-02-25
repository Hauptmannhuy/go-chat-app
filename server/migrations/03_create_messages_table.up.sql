CREATE TABLE IF NOT EXISTS messages (
 id SERIAL PRIMARY KEY,
 body VARCHAR(255) NOT NULL,
 user_id INTEGER REFERENCES users(id),
 message_id INTEGER NOT NULL, 
 chat_name VARCHAR(255) NOT NULL,
 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);