-- name: CreateFeedFollower :one

INSERT INTO feeds_follower (id,created_at,updated_at,user_id,feed_id)
VALUES ($1,$2,$3,$4,$5)
RETURNING *;

-- name: GetFeedFollower :many

SELECT * from feeds_follower WHERE user_id=$1;

-- name: DeleteFeedFollower :exec

DELETE FROM feeds_follower WHERE id=$1 and user_id=$2;

-- name: ExistFeedFollower :one
SELECT * FROM feeds_follower WHERE id=$1 and user_id=$2;