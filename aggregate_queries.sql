-- Count employees
SELECT COUNT(*) AS total_employees FROM employees;

-- Average salary
SELECT AVG(salary) AS avg_salary FROM employees;

-- Group by department
SELECT department, COUNT(*) AS total
FROM employees
GROUP BY department;