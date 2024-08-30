CREATE TABLE users
(
    id             SERIAL PRIMARY KEY,
    name           VARCHAR(100)       NOT NULL,
    email          VARCHAR(50) UNIQUE NOT NULL,
    username       VARCHAR(50) UNIQUE NOT NULL,
    password       VARCHAR(100)        NOT NULL,
    jwt_token      VARCHAR(100),
    deposit_amount DECIMAL(10, 2) DEFAULT 0,
    created_at     TIMESTAMP      DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP      DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE cars
(
    id           SERIAL PRIMARY KEY,
    model        VARCHAR(50)    NOT NULL,
    brand        VARCHAR(50)    NOT NULL,
    color        VARCHAR(50),
    category     VARCHAR(50),
    year         VARCHAR(50),
    rental_cost  DECIMAL(10, 2) NOT NULL,
    availability VARCHAR(50)
);

CREATE TABLE payments
(
    id          SERIAL PRIMARY KEY,
    type        VARCHAR(50),
    invoice_url VARCHAR(100),
    amount      DECIMAL(10, 2)
);

CREATE TABLE rental_histories
(
    id         SERIAL PRIMARY KEY,
    user_id    INT REFERENCES users (id),
    payment_id INT REFERENCES payments (id),
    start_date DATE,
    end_date   DATE
);

CREATE TABLE car_rentals
(
    id                SERIAL PRIMARY KEY,
    car_id            INT REFERENCES cars (id),
    rental_history_id INT REFERENCES rental_histories (id) ON DELETE CASCADE
);
