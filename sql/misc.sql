CREATE TABLE tasks(
    id INT PRIMARY KEY NOT NULL,
    name TEXT NOT NULL,
    idea TEXT NOT NULL,
    completed BOOLEAN
);


CREATE TABLE tags(
    id INT PRIMARY KEY NOT NULL,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE tag_task_relations(
    task_id INT,
    tag_id INT,
    PRIMARY KEY (task_id, tag_id),
    FOREIGN KEY (task_id) REFERENCES tasks(id),
    FOREIGN KEY (tag_id) REFERENCES tags(id)
);

SELECT 
    t.id, t.name, t.idea, t.completed,
    tg.id, tg.name
FROM 
    tasks t
JOIN 
    tag_task_relations rel ON t.id = rel.task_id
JOIN 
    tags tg ON rel.tag_id = tg.id;
