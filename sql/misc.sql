CREATE TABLE tags(
    DEDELtag_id INT,
    task_id INT,
)

CREATE TABLE tasks(
    id INT PRIMARY KEY NOT NULL,
    name TEXT NOT NULL,
    idea TEXT NOT NULL,
    completed BOOLEAN,

)


CREATE TABLE tags_tasks_relation(
    tag_id INT,
    task_id INT,
    PRIMARY KEY (tag_id, task_id),

    FOREIGN KEY (tag_id) REFERENCES tags(id),
    FOREIGN KEY (task_id) REFERENCES tasks(id)

);