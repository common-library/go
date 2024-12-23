-- migrate:up
CREATE TABLE IF NOT EXISTS test_01(
    field01 String,
	field02 UInt32
)
ENGINE = MergeTree
PRIMARY KEY (field01);

-- migrate:down
DROP TABLE IF EXISTS test_01;
