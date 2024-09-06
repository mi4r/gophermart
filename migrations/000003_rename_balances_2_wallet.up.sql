ALTER TABLE balances RENAME COLUMN value TO balance;
ALTER TABLE balances RENAME COLUMN bonuses TO scores;
ALTER TABLE balances RENAME TO wallets;

ALTER TABLE users RENAME COLUMN balance_id TO wallet_id;