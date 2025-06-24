-- RIGHT JOIN: Returns all records from the right table and the matched records from the left table
SELECT employees.first_name, departments.department_name
FROM employees
RIGHT JOIN departments ON employees.department_id = departments.department_id;