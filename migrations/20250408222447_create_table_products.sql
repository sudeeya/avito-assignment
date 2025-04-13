-- +goose Up
-- +goose StatementBegin
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reception_id UUID REFERENCES receptions(id) ON DELETE CASCADE NOT NULL,
    product_type_id UUID REFERENCES product_types(id) NOT NULL,
    datetime TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE products;
-- +goose StatementEnd
