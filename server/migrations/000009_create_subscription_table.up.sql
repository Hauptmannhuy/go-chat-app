CREATE TABLE IF NOT EXISTS subscriptions (
    id SERIAL PRIMARY KEY,               
    username VARCHAR(255) NOT NULL,  
    chatID VARCHAR(255) NOT NULL
);