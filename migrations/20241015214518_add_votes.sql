-- +goose Up
CREATE TABLE votes (
    chain TEXT NOT NULL,
    proposal_id TEXT NOT NULL,
    wallet TEXT NOT NULL,
    vote_option TEXT NOT NULL,
    vote_weight REAL NOT NULL,
    PRIMARY KEY (chain, proposal_id, wallet, vote_option)
);

-- +goose Down
DROP TABLE votes;
