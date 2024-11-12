-- +goose Up
CREATE TABLE proposals (
    chain TEXT NOT NULL,
    id TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    status TEXT NOT NULL,
    end_time TIMESTAMP NOT NULL,
    PRIMARY KEY (chain, id)
);

-- +goose Down
DROP TABLE proposals;
