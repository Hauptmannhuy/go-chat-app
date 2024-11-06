CREATE TABLE IF NOT EXISTS private_chats (
  id SERIAL PRIMARY KEY,
  chat_name VARCHAR(255) NOT NULL,
  user1_id INTEGER NOT NULL REFERENCES users(id),
  user2_id INTEGER NOT NULL REFERENCES users(id),
  UNIQUE(user1_id, user2_id)
)