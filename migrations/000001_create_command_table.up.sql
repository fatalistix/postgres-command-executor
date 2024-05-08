CREATE TABLE command (
    id BIGSERIAL UNIQUE NOT NULL,
    command TEXT UNIQUE NOT NULL,
    PRIMARY KEY(id)
);

CREATE INDEX command_command_idx ON command (command);