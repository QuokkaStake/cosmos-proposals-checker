-- +goose Up
CREATE TABLE mutes (
    chain TEXT,
    proposal_id TEXT,
    expires TIMESTAMP NOT NULL,
    comment TEXT,
    PRIMARY KEY (chain, proposal_id)
);

-- +goose Down
DROP TABLE mutes;
