-- +migrate Up
CREATE TABLE IF NOT EXISTS exams (
                                     id SERIAL PRIMARY KEY,
                                     instructor_id INTEGER NOT NULL,
                                     crn TEXT NOT NULL,
                                     date TEXT NOT NULL,
                                     created_at TIMESTAMP NOT NULL DEFAULT now(),
                                     updated_at TIMESTAMP NOT NULL DEFAULT now(),
                                     deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS answers (
                                       id SERIAL PRIMARY KEY,
                                       exam_id INTEGER NOT NULL REFERENCES exams(id) ON DELETE CASCADE,
                                       grade REAL,
                                       pdf_url TEXT NOT NULL,
                                       created_at TIMESTAMP NOT NULL DEFAULT now(),
                                       updated_at TIMESTAMP NOT NULL DEFAULT now(),
                                       deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS images (
                                      id SERIAL PRIMARY KEY,
                                      answer_id INTEGER NOT NULL REFERENCES answers(id) ON DELETE CASCADE,
                                      url TEXT NOT NULL,
                                      created_at TIMESTAMP NOT NULL DEFAULT now(),
                                      updated_at TIMESTAMP NOT NULL DEFAULT now(),
                                      deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS b_box_templates (
                                               id SERIAL PRIMARY KEY,
                                               exam_id INTEGER NOT NULL REFERENCES exams(id) ON DELETE CASCADE,
                                               name TEXT NOT NULL,
                                               created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                                               updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS b_box_meta_dbs (
                                              id SERIAL PRIMARY KEY,
                                              page INTEGER NOT NULL,
                                              template_id INTEGER NOT NULL REFERENCES b_box_templates(id) ON DELETE CASCADE,
                                              b_box_percent JSONB NOT NULL,
                                              created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                                              updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
