-- INNER JOIN: Returns records that have matching values in both tables
SELECT employees.first_name, departments.department_name
FROM employees
INNER JOIN departments ON employees.department_id = departments.department_id;