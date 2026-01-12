CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';


-- user table
CREATE TABLE user (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    login VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
);

CREATE TRIGGER update_user_updated_at
    BEFORE UPDATE ON user
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- balance table
CREATE TABLE balance (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    value BIGINT NOT NULL CHECK (value >= 0) DEFAULT 0,
    user_id UUID UNIQUE REFERENCES user (id) ON UPDATE CASCADE ON DELETE CASCADE,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
);

CREATE TRIGGER update_balance_updated_at
    BEFORE UPDATE ON balance
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- order table
CREATE TABLE order (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES user (id) ON UPDATE CASCADE ON DELETE CASCADE,
    order_num VARCHAR(255) NOT NULL UNIQUE,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
);

CREATE TRIGGER update_order_updated_at
    BEFORE UPDATE ON order
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- accrual table
CREATE TABLE accrual (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID UNIQUE REFERENCES order (id) ON UPDATE CASCADE ON DELETE CASCADE,
    status VARCHAR(10) NOT NULL CHECK (status IN ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED')),
    value BIGINT NOT NULL CHECK (value >= 0) DEFAULT 0,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
);

CREATE TRIGGER update_accrual_updated_at
    BEFORE UPDATE ON accrual
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- withdraw table
CREATE TABLE withdraw (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID UNIQUE REFERENCES order (id) ON UPDATE CASCADE ON DELETE CASCADE,
    value BIGINT NOT NULL CHECK (value >= 0) DEFAULT 0,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
);

CREATE TRIGGER update_withdraw_updated_at
    BEFORE UPDATE ON withdraw
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
