BEGIN;
CREATE TABLE IF NOT EXISTS epochs (
    epoch_number BIGINT NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    CONSTRAINT pk_epochs PRIMARY KEY(epoch_number)
);
CREATE TABLE IF NOT EXISTS slots (
    slot_number BIGINT NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    epoch_number BIGINT NOT NULL,
    CONSTRAINT pk_slots PRIMARY KEY(slot_number),
    CONSTRAINT fk_slots_epochs FOREIGN KEY(epoch_number) REFERENCES epochs(epoch_number) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS blocks (
    block_number BIGINT NOT NULL,
    block_root VARCHAR NOT NULL,
    state_root VARCHAR NOT NULL,
    slot_number BIGINT NOT NULL,
    gas_limit BIGINT NOT NULL,
    gas_used BIGINT NOT NULL,
    no_of_transactions INT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    CONSTRAINT pk_blocks PRIMARY KEY(block_number),
    CONSTRAINT fk_blocks_slots FOREIGN KEY(slot_number) REFERENCES slots(slot_number) ON DELETE CASCADE
);
COMMIT;