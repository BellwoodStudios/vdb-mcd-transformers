-- +goose Up
CREATE TABLE maker.median_kiss_batch
(
    id         SERIAL PRIMARY KEY,
    log_id     BIGINT     NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    header_id  INTEGER    NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    address_id INTEGER    NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    msg_sender INTEGER    NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    a_length   INTEGER    NOT NULL,
    a          TEXT ARRAY NOT NULL,
    UNIQUE (header_id, log_id)
);

CREATE INDEX median_kiss_batch_log_index
    ON maker.median_kiss_batch (log_id);
CREATE INDEX median_kiss_batch_header_index
    ON maker.median_kiss_batch (header_id);
CREATE INDEX median_kiss_batch_address_index
    ON maker.median_kiss_batch (address_id);
CREATE INDEX median_kiss_batch_msg_sender_index
    ON maker.median_kiss_batch (msg_sender);


-- +goose Down
DROP TABLE maker.median_kiss_batch;
