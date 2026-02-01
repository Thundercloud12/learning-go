-- +goose Up

ALTER TABLE url
ADD constraint unique_redirection UNIQUE(converted_url);
