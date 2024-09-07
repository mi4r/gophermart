CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    login VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    balance DECIMAL(10, 2) NOT NULL
);

CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    number VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL CHECK (status IN ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED')),
    accrual BIGINT NOT NULL,
    uploaded_at TIMESTAMP NOT NULL,
    user_login VARCHAR(255) NOT NULL,
    FOREIGN KEY (user_login) REFERENCES users(login)
);

