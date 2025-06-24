CREATE DATABASE IF NOT EXISTS testdb;
USE testdb;

CREATE TABLE IF NOT EXISTS users (
    id INT PRIMARY KEY,
    name VARCHAR(100),
    email VARCHAR(100)
);

INSERT INTO users (id, name, email) VALUES
(1, 'Surya', 'surya@example.com'),
(2, 'Alice', 'alice@example.com'),
(3, 'Bob', 'bob@example.com'),
(4, 'John', 'john@example.com')
(5,'craig','craig@example.com');
