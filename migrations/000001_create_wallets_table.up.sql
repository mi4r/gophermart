CREATE TABLE wallets (
    id SERIAL PRIMARY KEY,
    balance NUMERIC(15, 2) NOT NULL,
    scores INT NOT NULL,
    is_locked BOOLEAN DEFAULT FALSE
);