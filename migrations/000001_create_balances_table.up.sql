CREATE TABLE balances (
    id SERIAL PRIMARY KEY,
    value NUMERIC(15, 2) NOT NULL,
    bonuses INT NOT NULL
);