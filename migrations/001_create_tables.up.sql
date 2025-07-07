-- +migrate Up
CREATE TABLE IF NOT EXISTS classes (
                                       id SERIAL PRIMARY KEY,
                                       instructor_id INTEGER NOT NULL,
                                       crn TEXT NOT NULL,
                                       created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP
    );

CREATE TABLE IF NOT EXISTS exams (
                                     id SERIAL PRIMARY KEY,
                                     student_id INTEGER NOT NULL,
                                     class_id INTEGER NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    grade REAL,
    pdf_url TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP
    );

CREATE TABLE IF NOT EXISTS images (
                                      id SERIAL PRIMARY KEY,
                                      exam_id INTEGER NOT NULL REFERENCES exams(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP
    );