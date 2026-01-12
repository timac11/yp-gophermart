DROP TRIGGER IF EXISTS update_user_updated_at ON user;
DROP TRIGGER IF EXISTS update_balance_updated_at ON balance;
DROP TRIGGER IF EXISTS update_order_updated_at ON order;
DROP TRIGGER IF EXISTS update_accrual_updated_at ON accrual;
DROP TRIGGER IF EXISTS update_withdrawal_updated_at ON withdraw;

DROP TABLE IF EXISTS withdrawal;
DROP TABLE IF EXISTS accrual;
DROP TABLE IF EXISTS order;
DROP TABLE IF EXISTS balance;
DROP TABLE IF EXISTS user;

DROP FUNCTION IF EXISTS update_updated_at_column;
