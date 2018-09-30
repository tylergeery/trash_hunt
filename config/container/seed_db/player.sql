CREATE TABLE player (
    id SERIAL PRIMARY KEY,
    email VARCHAR(50),
    password VARCHAR(120),
    name VARCHAR(60),
    facebook_id VARCHAR(50) DEFAULT NULL,
    is_active SMALLINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'PST'),
    updated_at TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'PST')
);

INSERT INTO player (email, password, name)
VALUES ('tyler@geerydev.com', 'test', 'Tyler Geery');
