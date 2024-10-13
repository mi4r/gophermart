BEGIN;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status_enum') THEN
        CREATE TYPE status_enum AS ENUM ('NEW', 'REGISTERED', 'PROCESSING', 'INVALID', 'PROCESSED');
    END IF;
END
$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'reward_type_enum') THEN
        CREATE TYPE reward_type_enum AS ENUM ('%', 'pt');
    END IF;
END
$$;

CREATE TABLE rewards (
    id SERIAL PRIMARY KEY,
    match VARCHAR(255) UNIQUE NOT NULL,
    reward NUMERIC(10,2) NOT NULL,
    reward_type reward_type_enum DEFAULT '%' NOT NULL
);

CREATE TABLE goods (
    id SERIAL PRIMARY KEY,
    description VARCHAR(255) UNIQUE NOT NULL,
    price NUMERIC(10,2) NOT NULL
);

CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    order_number VARCHAR(255) UNIQUE NOT NULL,
    status status_enum DEFAULT 'REGISTERED' NOT NULL,
    accrual NUMERIC(10,2)
);

CREATE TABLE order_goods (
    order_id INT,
    good_id INT,
    PRIMARY KEY (order_id, good_id),
    FOREIGN KEY (order_id) REFERENCES orders(id),
    FOREIGN KEY (good_id) REFERENCES goods(id)
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    login VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(64) NOT NULL,
    current NUMERIC(10,2) DEFAULT 0 NOT NULL,
    withdrawn NUMERIC(10,2) DEFAULT 0 NOT NULL
);

CREATE TABLE user_orders (
    id SERIAL PRIMARY KEY,
    number VARCHAR(255) UNIQUE NOT NULL,
    status status_enum DEFAULT 'NEW' NOT NULL,
    sum NUMERIC(10,2) DEFAULT 0 NOT NULL,
    is_withdrawn BOOLEAN DEFAULT false NOT NULL,
    accrual NUMERIC(10,2) DEFAULT 0 NOT NULL,
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    processed_at TIMESTAMP,
    user_login VARCHAR(255) NOT NULL,
    FOREIGN KEY (user_login) REFERENCES users(login)
);

COMMIT;