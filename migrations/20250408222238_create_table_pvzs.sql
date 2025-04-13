-- +goose Up
-- +goose StatementBegin
CREATE TABLE pvzs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    registration_date TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    city_id UUID REFERENCES cities(id) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE pvzs;
-- +goose StatementEnd
