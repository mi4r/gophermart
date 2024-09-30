BEGIN;

CREATE TABLE rewards (
    id SERIAL PRIMARY KEY,
    match VARCHAR(255) UNIQUE NOT NULL,
    reward DOUBLE PRECISION NOT NULL,
    reward_type VARCHAR(2) DEFAULT '%' NOT NULL CHECK (reward_type IN ('%', 'pt'))
);

CREATE TABLE goods (
    id SERIAL PRIMARY KEY,
    description VARCHAR(255) UNIQUE NOT NULL,
    price DOUBLE PRECISION NOT NULL
);

CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    order_number VARCHAR(255) UNIQUE NOT NULL,
    status VARCHAR(50) DEFAULT 'NEW' NOT NULL CHECK (status IN ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED')),
    accrual DOUBLE PRECISION
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
    current DOUBLE PRECISION DEFAULT 0 NOT NULL,
    withdrawn DOUBLE PRECISION DEFAULT 0 NOT NULL
);

CREATE TABLE user_orders (
    id SERIAL PRIMARY KEY,
    number VARCHAR(255) UNIQUE NOT NULL,
    status VARCHAR(50) DEFAULT 'NEW' NOT NULL CHECK (status IN ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED')),
    sum DOUBLE PRECISION DEFAULT 0 NOT NULL,
    is_withdrawn BOOLEAN DEFAULT false NOT NULL,
    accrual DOUBLE PRECISION DEFAULT 0 NOT NULL,
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    processed_at TIMESTAMP,
    user_login VARCHAR(255) NOT NULL,
    FOREIGN KEY (user_login) REFERENCES users(login)
);

COMMIT;