-- name: CreateTable01 :exec
CREATE TABLE IF NOT EXISTS table01_for_mysql (
  field01 VARCHAR(64) PRIMARY KEY,
  field02 BIGINT NOT NULL
);

-- name: DropTable01 :exec
DROP TABLE IF EXISTS table01_for_mysql;

-- name: GetTable01 :one
SELECT * FROM table01_for_mysql WHERE field01 = ?;

-- name: ListTable01 :many
SELECT * FROM table01_for_mysql ORDER BY field01;

-- name: InsertTable01 :execresult
INSERT INTO table01_for_mysql(field01, field02) VALUES(?, ?);

-- name: UpdateTable01 :exec
UPDATE table01_for_mysql set field02 = ? WHERE field01 = ?;

-- name: DeleteTable01 :exec
DELETE FROM table01_for_mysql WHERE field01 = ?;
