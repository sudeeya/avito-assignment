-- +goose Up
-- +goose StatementBegin
CREATE TABLE product_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT UNIQUE NOT NULL
);

INSERT INTO product_types (name)
VALUES
    ('электроника'),
    ('одежда'),
    ('обувь');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE product_types;
-- +goose StatementEnd
