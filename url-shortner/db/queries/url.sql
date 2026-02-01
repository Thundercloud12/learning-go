-- name: CreateMapping :exec
INSERT INTO url (id, original_url, converted_url)
VALUES ($1, $2, $3);

-- name: GetConverted :one
SELECT * FROM url WHERE converted_url = $1;

-- name: GetAllMappings :many
SELECT * FROM url;

-- name: ExistsUrl :one
SELECT EXISTS(
    SELECT 1
    FROM url
    WHERE original_url=$1
)AS url_exists;