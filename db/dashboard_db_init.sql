PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE asset (
	symbol TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	type TEXT NOT NULL, -- 'forex', 'crypto', 'stock', 'commodity', 'bond'
	value_usd REAL,
	last_synched TIMESTAMP
);
CREATE TABLE wallet(
	name TEXT PRIMARY KEY,
	value_usd REAL DEFAULT 0
);
CREATE TABLE transfer(
	id INTEGER PRIMARY KEY,
	timestamp_utc TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	ammount REAL NOT NULL,
	asset_symbol TEXT NOT NULL,
	from_wallet_id INTEGER,
	to_wallet_id INTEGER NOT NULL,
	FOREIGN KEY(asset_symbol) REFERENCES asset(symbol)
);
CREATE TABLE price_daily(
	id INTEGER PRIMARY KEY,
	asset_symbol TEXT NOT NULL,
	timestamp_utc TIMESTAMP NOT NULL,
	value_usd REAL NOT NULL,
	FOREIGN KEY(asset_symbol) REFERENCES asset(symbol)
);
INSERT INTO asset(symbol, name, type) VALUES('USD', 'United States Dollar', 'forex');
INSERT INTO asset(symbol, name, type) VALUES('EUR', 'Euro', 'forex');
INSERT INTO asset(symbol, name, type) VALUES('BTC', 'Bitcoin', 'crypto');
COMMIT;
