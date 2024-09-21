-- name: CreateTable01 :exec
CREATE TABLE IF NOT EXISTS table01_for_postgresql (
  field01 TEXT PRIMARY KEY,
  field02 BIGINT NOT NULL
);

-- name: DropTable01 :exec
DROP TABLE IF EXISTS table01_for_postgresql;

-- name: GetTable01 :one
SELECT * FROM table01_for_postgresql WHERE field01 = $1;

-- name: ListTable01 :many
SELECT * FROM table01_for_postgresql ORDER BY field01;

-- name: InsertTable01 :one
INSERT INTO table01_for_postgresql(field01, field02) VALUES($1, $2) RETURNING *;

-- name: UpdateTable01 :exec
UPDATE table01_for_postgresql set field02 = $2 WHERE field01 = $1;

-- name: DeleteTable01 :exec
DELETE FROM table01_for_postgresql WHERE field01 = $1;
