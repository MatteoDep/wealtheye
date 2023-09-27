PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE asset (name text primary key , symbol text);
INSERT INTO asset VALUES('usd','USD');
INSERT INTO asset VALUES('eur','EUR');
INSERT INTO asset VALUES('btc','BTC');
COMMIT;
