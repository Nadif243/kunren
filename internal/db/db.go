package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

// InitDB establishes the crucible and constructs the structural geometry (tables).
func InitDB(filepath string) {
	var err error
	DB, err = sql.Open("sqlite", filepath)
	if err != nil {
		log.Fatal("Critical: Failed to ignite the database crucible:", err)
	}

	createTables := `
	CREATE TABLE IF NOT EXISTS users (
		discord_id TEXT PRIMARY KEY,
		bound_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS vaults (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		owner_id TEXT,
		name TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (owner_id) REFERENCES users(discord_id)
	);

	CREATE TABLE IF NOT EXISTS sentences (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		vault_id INTEGER,
		keyword TEXT,
		raw_text TEXT,
		source_url TEXT,
		saved_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (vault_id) REFERENCES vaults(id)
	);`

	_, err = DB.Exec(createTables)
	if err != nil {
		log.Fatal("Critical: Failed to engrave the schema:", err)
	}
}

// BindUser ties a Discord soul to the engine and creates their first grimoire.
func BindUser(discordID string) (string, error) {
	// Check if already bound
	var exists bool
	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE discord_id = ?)", discordID).Scan(&exists)
	if err != nil {
		return "", err
	}
	if exists {
		return "Your soul is already bound to the crucible. Use `!vault create <name>` to expand your grimoires.", nil
	}

	// Begin the binding transaction
	tx, err := DB.Begin()
	if err != nil {
		return "", err
	}

	_, err = tx.Exec("INSERT INTO users (discord_id) VALUES (?)", discordID)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	_, err = tx.Exec("INSERT INTO vaults (owner_id, name) VALUES (?, 'Primary Grimoire')", discordID)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	err = tx.Commit()
	if err != nil {
		return "", err
	}

	return "Binding complete. A 'Primary Grimoire' has been forged in the void.", nil
}

// CreateVault spins up a new isolated collection for the user.
func CreateVault(discordID, vaultName string) (string, error) {
	var exists bool
	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE discord_id = ?)", discordID).Scan(&exists)
	if err != nil {
		return "", err
	}
	if !exists {
		return "You are not bound to the crucible. Run `!bind` first.", nil
	}

	_, err = DB.Exec("INSERT INTO vaults (owner_id, name) VALUES (?, ?)", discordID, vaultName)
	if err != nil {
		return "", fmt.Errorf("failed to forge the vault: %w", err)
	}

	return fmt.Sprintf("Vault 『%s』 has been successfully forged.", vaultName), nil
}
