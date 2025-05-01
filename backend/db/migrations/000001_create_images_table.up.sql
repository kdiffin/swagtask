-- 0001_create_images_table.up.sql

CREATE TABLE task_images (
  id SERIAL PRIMARY KEY,
  task_id INTEGER REFERENCES tasks(id) ON DELETE CASCADE,
  file_path TEXT NOT NULL, -- e.g. "uploads/img_1234.png"
  uploaded_at TIMESTAMP DEFAULT now()
);
