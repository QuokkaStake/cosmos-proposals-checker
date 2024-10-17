-- +goose Up
CREATE TABLE query_last_block (
    chain TEXT NOT NULL,
    query TEXT NOT NULL,
    height INTEGER NOT NULL,
    PRIMARY KEY (chain, query)
);

-- +goose Down
DROP TABLE query_last_block;
