CREATE TABLE IF NOT EXISTS projects (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    user_id INT NOT NULL
);

CREATE TABLE IF NOT EXISTS todo_items (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    project_id INT NOT NULL REFERENCES projects
);