-- +goose Up
-- +goose StatementBegin
CREATE TYPE reception_status AS ENUM ('in_progress', 'close');

CREATE TABLE receptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pvz_id UUID REFERENCES pvzs(id) ON DELETE CASCADE NOT NULL,
    datetime TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    status reception_status DEFAULT 'in_progress' NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE receptions;

DROP TYPE reception_status;
-- +goose StatementEnd
