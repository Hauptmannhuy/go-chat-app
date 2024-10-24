CREATE TABLE IF NOT EXISTS subscriptions (
    id SERIAL PRIMARY KEY,               
    username VARCHAR(255) NOT NULL,  
    chat_id VARCHAR(255) NOT NULL
);