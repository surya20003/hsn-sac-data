-- LEFT JOIN: Returns all records from the left table and the matched records from the right table
SELECT employees.first_name, departments.department_name
FROM employees
LEFT JOIN departments ON employees.department_id = departments.department_id;