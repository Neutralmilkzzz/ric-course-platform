-- PostgreSQL schema for RIC 选课平台
CREATE TABLE IF NOT EXISTS students (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS courses (
    id SERIAL PRIMARY KEY,
    code TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS enrollments (
    student_id INT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    course_id  INT NOT NULL REFERENCES courses(id)  ON DELETE CASCADE,
    PRIMARY KEY (student_id, course_id)
);
