-- FULL OUTER JOIN: Returns all records when there is a match in either left or right table
-- Note: Not supported in MySQL directly. Use UNION of LEFT and RIGHT JOIN
SELECT employees.first_name, departments.department_name
FROM employees
LEFT JOIN departments ON employees.department_id = departments.department_id
UNION
SELECT employees.first_name, departments.department_name
FROM employees
RIGHT JOIN departments ON employees.department_id = departments.department_id;