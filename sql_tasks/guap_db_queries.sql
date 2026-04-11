CREATE TABLE students (
    id          SERIAL PRIMARY KEY,
    full_name   VARCHAR(150) NOT NULL,
    group_name  VARCHAR(20),
    email       VARCHAR(255) UNIQUE,
    enrolled_at DATE,
    is_active   BOOLEAN DEFAULT TRUE
);

CREATE TABLE subjects (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(200) NOT NULL,
    department  VARCHAR(100),
    credits     INTEGER
);

CREATE TABLE grades (
    id          SERIAL PRIMARY KEY,
    student_id  INTEGER REFERENCES students(id),
    subject_id  INTEGER REFERENCES subjects(id),
    grade       INTEGER CHECK (grade BETWEEN 2 AND 5),
    exam_date   DATE,
    semester    INTEGER
);

CREATE TABLE schedule (
    id          SERIAL PRIMARY KEY,
    subject_id  INTEGER REFERENCES subjects(id),
    group_name  VARCHAR(20),
    weekday     VARCHAR(20),
    time_start  TIME,
    room        VARCHAR(20)
);

SELECT id, full_name, email, enrolled_at
FROM students
WHERE group_name = 'Z3420'
  AND is_active = TRUE
ORDER BY full_name;

SELECT group_name, COUNT(*) AS students_count
FROM students
WHERE is_active = TRUE
GROUP BY group_name
ORDER BY students_count DESC;

SELECT s.id,
       s.full_name,
       s.group_name,
       ROUND(AVG(g.grade), 2) AS avg_grade,
       COUNT(g.id)            AS exams_passed
FROM students s
JOIN grades g ON g.student_id = s.id
WHERE g.semester = 3
GROUP BY s.id, s.full_name, s.group_name
ORDER BY avg_grade DESC;

SELECT s.full_name,
       s.group_name,
       sub.name AS subject,
       g.grade,
       g.exam_date
FROM grades g
JOIN students s  ON s.id = g.student_id
JOIN subjects sub ON sub.id = g.subject_id
WHERE g.grade = 2
ORDER BY g.exam_date DESC;

SELECT email, COUNT(*) AS count
FROM students
GROUP BY email
HAVING COUNT(*) > 1;

SELECT sch.time_start,
       sub.name    AS subject,
       sch.room,
       sch.weekday
FROM schedule sch
JOIN subjects sub ON sub.id = sch.subject_id
WHERE sch.group_name = 'Z3420'
  AND sch.weekday = 'Понедельник'
ORDER BY sch.time_start;

SELECT s.id, s.full_name, s.group_name
FROM students s
LEFT JOIN grades g ON g.student_id = s.id
WHERE g.id IS NULL
  AND s.is_active = TRUE;

SELECT sub.id, sub.name, sub.department
FROM subjects sub
LEFT JOIN grades g ON g.subject_id = sub.id
WHERE g.id IS NULL;

SELECT s.full_name,
       s.group_name,
       ROUND(AVG(g.grade), 2) AS avg_grade
FROM students s
JOIN grades g ON g.student_id = s.id
GROUP BY s.id, s.full_name, s.group_name
ORDER BY avg_grade DESC
LIMIT 5;

SELECT TO_CHAR(exam_date, 'Day') AS weekday,
       COUNT(*)                  AS exams_count
FROM grades
GROUP BY TO_CHAR(exam_date, 'Day')
ORDER BY exams_count DESC;

SELECT 'students.full_name'  AS field, COUNT(*) AS nulls FROM students WHERE full_name IS NULL
UNION ALL
SELECT 'students.group_name', COUNT(*) FROM students WHERE group_name IS NULL
UNION ALL
SELECT 'grades.grade',        COUNT(*) FROM grades WHERE grade IS NULL
UNION ALL
SELECT 'grades.student_id',   COUNT(*) FROM grades WHERE student_id IS NULL;

DELETE FROM grades
WHERE student_id IN (
    SELECT id FROM students WHERE email LIKE '%@test.guap.ru'
);

DELETE FROM students
WHERE email LIKE '%@test.guap.ru';
