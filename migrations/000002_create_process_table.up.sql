CREATE TYPE process_status AS ENUM('executing', 'finished', 'error');

CREATE TABLE process (
    id UUID UNIQUE NOT NULL,
    output TEXT NOT NULL DEFAULT '',
    error TEXT NOT NULL DEFAULT '',
    status process_status NOT NULL DEFAULT 'executing',
    exit_code INT NOT NULL DEFAULT -1,
    PRIMARY KEY(id)
);