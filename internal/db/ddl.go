package db

import "SB/logger"

func InitMigrations() error {
	usersTableQuery := `
		CREATE TABLE IF NOT EXISTS users(
    id SERIAL PRIMARY KEY,
    full_name VARCHAR NOT NULL,
    password VARCHAR NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    active BOOLEAN DEFAULT true,
    deleted_at TIMESTAMP
);`
	_, err := db.Exec(usersTableQuery)
	if err != nil {
		logger.Error.Printf("[db] InitMigrations(): error during create users table: %v", err.Error())
		return err
	}

	accountsTableQuery := `
		CREATE TABLE IF NOT EXISTS accounts(
	id SERIAL PRIMARY KEY,
	user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	phone_number VARCHAR NOT NULL,
	balance BIGINT NOT NULL DEFAULT 0,
	currency VARCHAR(3) NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP,
	deleted_at TIMESTAMP,
    active BOOLEAN DEFAULT TRUE
);`

	_, err = db.Exec(accountsTableQuery)
	if err != nil {
		logger.Error.Printf("[db] InitMigrations(): error during create accounts table: %v", err.Error())
		return err
	}

	transactionsTableQuery := `
		CREATE TABLE IF NOT EXISTS transactions (
	id SERIAL PRIMARY KEY,
	from_account_id INT REFERENCES accounts(id),
	to_account_id INT REFERENCES accounts(id),
	amount BIGINT NOT NULL CHECK (amount > 0),
	currency char(3) not null, 
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`

	_, err = db.Exec(transactionsTableQuery)
	if err != nil {
		logger.Error.Printf("[db] InitMigrations(): error during create transactions table: %v", err.Error())
		return err
	}

	creditsTableQuery := `
		CREATE TABLE IF NOT EXISTS credits (
	id SERIAL PRIMARY KEY,
	user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	amount BIGINT NOT NULL CHECK (amount > 0),
    currency VARCHAR(3) NOT NULL,
	duration_months INT NOT NULL CHECK (duration_months > 0),
	interest_rate NUMERIC(5,2) NOT NULL CHECK (interest_rate >= 0),
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	approved_at TIMESTAMP,
    active BOOLEAN
);`

	_, err = db.Exec(creditsTableQuery)
	if err != nil {
		logger.Error.Printf("[db] InitMigrations(): error during create credits table: %v", err.Error())
		return err
	}
	depositsTableQuery := `
		CREATE TABLE IF NOT EXISTS deposits (
	id SERIAL PRIMARY KEY,
	user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	amount BIGINT NOT NULL CHECK (amount > 0),
	currency VARCHAR(3) NOT NULL,
	interest_rate NUMERIC(5,2) NOT NULL CHECK (interest_rate >= 0),
	duration_months INT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	expires_at TIMESTAMP NOT NULL,
    active BOOLEAN
);`

	_, err = db.Exec(depositsTableQuery)
	if err != nil {
		logger.Error.Printf("[db] InitMigrations(): error during create deposits table: %v", err.Error())
		return err
	}

	auditlogsQuery := `
		CREATE TABLE IF NOT EXISTS audit_logs (
	id SERIAL PRIMARY KEY,
	action TEXT,
	entity TEXT,
	entity_id INT,
	user_id INT,
	timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`
	_, err = db.Exec(auditlogsQuery)
	if err != nil {
		logger.Error.Printf("[db] InitMigrations(): error during create audit_logs table: %v", err.Error())
		return err
	}

	entryQuery := `
		CREATE TABLE IF NOT EXISTS entries (
		    id bigserial primary key, 
		    account_id bigint not null,
		    amount bigint not null,
		    created_at timestamp default current_timestamp
		)`

	_, err = db.Exec(entryQuery)
	if err != nil {
		logger.Error.Printf("[db] InitMigrations(): error during create entries table: %v", err.Error())
		return err
	}

	return nil
}
