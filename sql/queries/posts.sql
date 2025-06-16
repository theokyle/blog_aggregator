-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    $3,
    $4,
    $5,
    $6
) RETURNING *;

-- name: GetPostsForUser :many
SELECT * FROM posts 
JOIN feeds ON posts.feed_id = feeds.id 
JOIN feed_follows ON feeds.id = feed_follows.feed_id
WHERE feed_follows.user_id = $1 
ORDER BY published_at DESC LIMIT $2;