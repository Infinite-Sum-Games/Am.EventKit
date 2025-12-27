-- name: GetAllDisputesQuery :many
SELECT id,
    txn_id,
    student_email,
    description,
    event_id,
    dispute_status
FROM dispute;

-- name: GetDisputeByIDQuery :one
SELECT id,
    txn_id,
    student_email,
    description,
    event_id,
    dispute_status
FROM dispute
WHERE id = $1;

-- name: CreateDisputeQuery :exec
INSERT INTO dispute (
    txn_id,
    event_id
) VALUES ($1, $2);

-- name: GetEventIdByTxnIdQuery :one
SELECT 
    b.event_id AS event_id,
    e.event_status AS event_status
FROM bookings b
INNER JOIN event e on b.event_id = e.id
WHERE txn_id = $1;

-- name: IncrementSeatFilledCountQuery :execrows
UPDATE event
SET seats_filled = seats_filled + 1
WHERE id = $1;

-- name: GetEventIdByDisputeIDQuery :one
SELECT 
    d.event_id AS event_id
FROM dispute d
WHERE d.id = $1;

-- name: DecrementSeatFilledCountQuery :execrows
UPDATE event
SET seats_filled = seats_filled - 1
WHERE id = $1;

-- name: UpdateDisputeQuery :execrows
UPDATE dispute
SET student_email = $2,
    description = $3,
    updated_at = NOW()
WHERE id = $1;

-- name: CloseAsTrueDisputeQuery :execrows
UPDATE dispute
SET dispute_status = 'CLOSED_AS_TRUE',
    updated_at = NOW()
WHERE id = $1;

-- name: CloseAsFalseDisputeQuery :execrows
UPDATE dispute
SET dispute_status = 'CLOSED_AS_FALSE',
    updated_at = NOW()
WHERE id = $1;


