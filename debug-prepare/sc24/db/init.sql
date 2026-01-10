CREATE TABLE IF NOT EXISTS workers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    script_content TEXT NOT NULL,
    cpu_limit INT DEFAULT 10,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO workers (name, script_content, cpu_limit) VALUES ('router-worker', 'console.log("init")', 50);