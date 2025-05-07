-- Create sequences
CREATE SEQUENCE tag_task_relations_id_seq START 1 INCREMENT 1;
CREATE SEQUENCE tags_id_seq START 1 INCREMENT 1;
CREATE SEQUENCE tasks_id_seq START 1 INCREMENT 1;

-- Create tables
CREATE TABLE tags (
    id INTEGER PRIMARY KEY DEFAULT nextval('tags_id_seq'),
    name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()

);

CREATE TABLE tasks (
    id INTEGER PRIMARY KEY DEFAULT nextval('tasks_id_seq'),
    name TEXT NOT NULL,
    idea TEXT NOT NULL UNIQUE,
    completed BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE tag_task_relations (
    id INTEGER PRIMARY KEY DEFAULT nextval('tag_task_relations_id_seq'),
    tag_id INTEGER NOT NULL,
    task_id INTEGER NOT NULL,
    CONSTRAINT tag_task_unique UNIQUE (tag_id, task_id),
    CONSTRAINT fk_tag FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE,
    CONSTRAINT fk_task FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
);
