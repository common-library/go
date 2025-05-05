-- migrate:up
ALTER TABLE test_01 ADD COLUMN field03 String;

-- migrate:down
ALTER TABLE test_01 DROP COLUMN field03;
