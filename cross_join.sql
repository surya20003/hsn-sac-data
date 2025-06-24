-- CROSS JOIN: Returns the Cartesian product of both tables
SELECT employees.first_name, departments.department_name
FROM employees
CROSS JOIN departments;