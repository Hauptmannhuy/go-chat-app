CREATE TABLE IF NOT EXISTS group_chats (
   id SERIAL PRIMARY KEY,
   creator_id int NOT NULL REFERENCES users(id), 
   chat_name VARCHAR(255)  NOT NULL

);