-- SELF JOIN: Join a table with itself
SELECT A.first_name AS Employee, B.first_name AS Manager
FROM employees A
INNER JOIN employees B ON A.manager_id = B.employee_id;