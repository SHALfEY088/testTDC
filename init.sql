CREATE TABLE IF NOT EXISTS clients (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    balance INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    client_id INTEGER NOT NULL,
    amount INTEGER NOT NULL,
    result VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

INSERT INTO clients (id, name, balance)
VALUES (0, 'Andrew', 999999),
       (1, 'Bill', 88888),
       (2, 'Clarc',7777),
       (3, 'David', 555),
       (4, 'Egor', 44),
       (5, 'Franklin', 2)
