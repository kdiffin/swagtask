-- Insert sample tags
INSERT INTO tags (id, name) VALUES
(1, 'urgent'),
(2, 'home'),
(3, 'work'),
(4, 'long-term'),
(5, 'school');

-- Insert sample tasks
INSERT INTO tasks (id, name, idea, completed) VALUES
(1, 'Buy groceries', 'Get milk, eggs, and bread', false),
(2, 'Finish CS homework', 'Implement Dijkstra’s algorithm', false),
(3, 'Clean the garage', 'Organize tools and sweep the floor', true),
(4, 'Prepare presentation', 'Slides for Monday’s team meeting', false),
(5, 'Read networking notes', 'Review TCP/IP before quiz', true);

-- Link tags to tasks
INSERT INTO tag_task_relations (task_id, tag_id) VALUES
(1, 1), -- Buy groceries - urgent
(1, 2), -- Buy groceries - home
(2, 3), -- CS homework - work
(2, 5), -- CS homework - school
(3, 2), -- Clean garage - home
(4, 3), -- Presentation - work
(4, 1), -- Presentation - urgent
(5, 4), -- Read notes - long-term
(5, 5); -- Read notes - school


