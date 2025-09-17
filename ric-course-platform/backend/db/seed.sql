-- Sample data
INSERT INTO students (name) VALUES
 ('Alice'), ('Bob'), ('Charlie'), ('Diana'), ('Ethan')
ON CONFLICT DO NOTHING;

INSERT INTO courses (code, title) VALUES
 ('CS101', 'Intro to Computer Science'),
 ('CS102', 'Data Structures'),
 ('CS103', 'Algorithms'),
 ('MATH101', 'Calculus I'),
 ('MATH102', 'Linear Algebra'),
 ('STAT101', 'Statistics'),
 ('DS101', 'Intro to Data Science'),
 ('AI101', 'Intro to Artificial Intelligence')
ON CONFLICT DO NOTHING;

-- Enrollments (many-to-many)
INSERT INTO enrollments (student_id, course_id)
SELECT s.id, c.id
FROM students s, courses c
WHERE (s.name, c.code) IN (
 ('Alice','CS101'),
 ('Alice','MATH101'),
 ('Bob','CS101'),
 ('Bob','CS102'),
 ('Charlie','CS103'),
 ('Charlie','MATH102'),
 ('Diana','AI101'),
 ('Diana','DS101'),
 ('Ethan','STAT101'),
 ('Ethan','MATH101')
)
ON CONFLICT DO NOTHING;
