-- name: CreateMapping :exec
INSERT INTO url (id, original_url, converted_url)
VALUES ($1, $2, $3);

-- name: GetConverted :one
SELECT * FROM url WHERE converted_url = $1;