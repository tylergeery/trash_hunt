CREATE TABLE game (
    id SERIAL PRIMARY KEY,
    player1_id integer NOT NULL,
    player2_id integer NOT NULL,
    winner_id integer DEFAULT NULL,
    loser_id DEFAULT NULL,
    status smallint DEFAULT 1,
    ended_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'PST'),
    updated_at TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'PST')
);
