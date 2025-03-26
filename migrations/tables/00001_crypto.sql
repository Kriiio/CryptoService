-- +goose Up
-- +goose StatementBegin
CREATE TABLE crypto (
    id SERIAL PRIMARY KEY,
    time TIMESTAMP,
    ask_price FLOAT,
    ask_volume FLOAT,
    ask_time TIMESTAMP,
    bid_price FLOAT,
    bid_volume FLOAT,
    bid_time TIMESTAMP
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE crypto;
-- +goose StatementEnd
