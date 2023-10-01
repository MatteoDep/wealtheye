PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE asset (
	symbol TEXT PRIMARY KEY,
	value_usd REAL,
	last_synched TEXT
);
CREATE TABLE wallet(
	id INTEGER PRIMARY KEY,
	label TEXT NOT NULL,
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
	timestamp_utc TEXT NOT NULL,
	value_usd REAL NOT NULL,
	FOREIGN KEY(asset_symbol) REFERENCES asset(symbol)
);
INSERT INTO asset(symbol) VALUES('USD');
INSERT INTO asset(symbol) VALUES('EUR');
INSERT INTO asset(symbol) VALUES('BTC');
COMMIT;
