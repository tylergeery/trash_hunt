CREATE TABLE test (
    id SERIAL PRIMARY KEY,
    keyword VARCHAR(20),
    created_at TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'PST'),
    updated_at TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'PST')
);

COMMENT ON COLUMN test.keyword IS 'Testing 1,2,3';
