-- migrate:up
CREATE TABLE IF NOT EXISTS test_01(
    field01 VARCHAR(10) NOT NULL,
    field02 INTEGER NOT NULL
);

-- migrate:down
DROP TABLE IF EXISTS test_01;
