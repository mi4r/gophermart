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
    status VARCHAR(50) DEFAULT 'NEW' NOT NULL CHECK (status IN ('NEW', 'REGISTERED', 'PROCESSING', 'INVALID', 'PROCESSED')),
    accrual DOUBLE PRECISION
);

CREATE TABLE order_goods (
    order_id INT,
    good_id INT,
    PRIMARY KEY (order_id, good_id),
    FOREIGN KEY (order_id) REFERENCES orders(id),
    FOREIGN KEY (good_id) REFERENCES goods(id)
);

COMMIT;