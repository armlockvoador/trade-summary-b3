-- name: CreateTrade :exec
INSERT INTO trade (
    id,
    close_time,
    trade_date,
    instrument_code,
    trade_price,
    trade_quantity,
    created_at,
    updated_at,
    deleted
) VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8, false
);