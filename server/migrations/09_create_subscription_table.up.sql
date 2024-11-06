CREATE TABLE IF NOT EXISTS group_chat_subs (
    id SERIAL PRIMARY KEY,               
    user_id INTEGER NOT NULL REFERENCES users(id),
    chat_id INTEGER NOT NULL REFERENCES group_chats(id)
);