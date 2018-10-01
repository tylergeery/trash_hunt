CREATE TABLE player (
    id SERIAL PRIMARY KEY,
    email VARCHAR(50),
    password VARCHAR(120),
    name VARCHAR(60),
    facebook_id VARCHAR(50) DEFAULT NULL,
    status VARCHAR(50) DEFAULT 'Active',
    created_at TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'PST'),
    updated_at TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'PST')
);

INSERT INTO player (email, password, name)
VALUES ('tyler@geerydev.com', 'test', 'Tyler Geery');
