-- name: GetAllStudents :many
SELECT 
  name AS student_name, 
  email, 
  phone_number, 
  college_name AS college
FROM student 
ORDER BY student_name ASC;
