-- +goose Up

CREATE TABLE url (
    id UUID PRIMARY KEY,
    original_url TEXT NOT NULL,
    converted_url TEXT NOT NULL
);

-- +goose Down

DROP TABLE IF EXISTS url;
