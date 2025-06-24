-- Select all columns
SELECT * FROM employees;

-- Select specific columns
SELECT first_name, last_name FROM employees;

-- With WHERE clause
SELECT * FROM employees WHERE department = 'HR';

-- With ORDER BY
SELECT * FROM employees ORDER BY salary DESC;