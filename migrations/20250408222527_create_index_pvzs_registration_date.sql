-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_receptions_status_datetime ON receptions(status, datetime);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_receptions_datetime;
-- +goose StatementEnd
