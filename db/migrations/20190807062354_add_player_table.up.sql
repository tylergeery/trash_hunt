CREATE TABLE player (
    id SERIAL PRIMARY KEY,
    email VARCHAR(120) UNIQUE,
    password VARCHAR(200),
    username VARCHAR(50) UNIQUE,
    token VARCHAR(200),
    status smallint DEFAULT 1,
    created_at TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'PST'),
    updated_at TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'PST')
);
