-- Up Migration (create table)
CREATE TABLE orders (
    id VARCHAR(255) PRIMARY KEY,
    price DECIMAL(10, 2) NOT NULL,
    tax DECIMAL(10, 2) NOT NULL,
    final_price DECIMAL(10, 2) NOT NULL
);