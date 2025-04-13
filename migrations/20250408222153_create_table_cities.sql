-- +goose Up
-- +goose StatementBegin
CREATE TABLE cities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT UNIQUE NOT NULL
);

INSERT INTO cities (name)
VALUES
    ('Москва'),
    ('Санкт-Петербург'),
    ('Казань');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE cities;
-- +goose StatementEnd
