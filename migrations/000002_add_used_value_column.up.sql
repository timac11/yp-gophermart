ALTER TABLE balance ADD used_value BIGINT NOT NULL CHECK (used_value >= 0) DEFAULT 0;
