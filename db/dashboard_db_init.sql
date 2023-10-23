PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE asset (
	symbol TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	type TEXT NOT NULL -- 'forex', 'crypto', 'stock', 'commodity', 'bond'
);
CREATE TABLE wallet(
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL,
	value_usd REAL NOT NULL DEFAULT 0
);
CREATE TABLE transfer(
	id INTEGER PRIMARY KEY,
	timestamp_utc TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	ammount REAL NOT NULL,
	asset_symbol TEXT NOT NULL,
	from_wallet_id INTEGER,
	to_wallet_id INTEGER,
	FOREIGN KEY(asset_symbol) REFERENCES asset(symbol),
	FOREIGN KEY(from_wallet_id) REFERENCES wallet(id),
	FOREIGN KEY(to_wallet_id) REFERENCES wallet(id)
);
CREATE TABLE price_daily(
	id INTEGER PRIMARY KEY,
	asset_symbol TEXT NOT NULL,
	timestamp_utc TIMESTAMP NOT NULL,
	value_usd REAL NOT NULL,
	FOREIGN KEY(asset_symbol) REFERENCES asset(symbol),
	UNIQUE(asset_symbol, timestamp_utc)
);
INSERT INTO asset(symbol, name, type) VALUES('USD', 'United States Dollar', 'forex');
INSERT INTO asset(symbol, name, type) VALUES('EUR', 'Euro', 'forex');
INSERT INTO asset(symbol, name, type) VALUES('BTC', 'Bitcoin', 'crypto');
COMMIT;
