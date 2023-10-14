CREATE TYPE transaction_account_sign_type AS ENUM ('CREDIT','DEBIT');

CREATE TABLE transactions_accounttype (
    accounttype_id integer NOT NULL PRIMARY KEY,
    accounttype_name varchar(20) NOT NULL DEFAULT '',
    accounttype_sign transaction_account_sign_type NOT NULL DEFAULT 'CREDIT'
);

INSERT INTO transactions_accounttype VALUES (1,'ASSET','DEBIT'),(2,'LIABILITY','CREDIT'),
                                            (3,'EQUITY','CREDIT'),(4,'INCOME','CREDIT'),
                                            (5,'EXPENSE','DEBIT'),(6,'GAIN','CREDIT'),(7,'LOSS','DEBIT');

CREATE TABLE transactions_accounts (
    account_id integer NOT NULL PRIMARY KEY,
    account_name varchar(50) DEFAULT '',
    accounttype_id integer NOT NULL DEFAULT '0',
    account_memo varchar(70) DEFAULT '',
    account_starting decimal(12,2) DEFAULT '0.00',
    account_current bool NOT NULL DEFAULT true,
    account_left integer DEFAULT NULL,
    account_right integer DEFAULT NULL,
    account_balance decimal(12,2) DEFAULT NULL,
    account_subtotal decimal(12,2) DEFAULT NULL,
    account_fullname varchar(200) DEFAULT '',
    account_parent integer NOT NULL DEFAULT '0',
    account_reconcile_date date DEFAULT NULL,
    account_flagged bool NOT NULL DEFAULT false,
    account_locked bool NOT NULL DEFAULT false,
    account_open_date timestamp without time zone DEFAULT now(),
    account_close_date timestamp without time zone DEFAULT NULL,
    account_code varchar(50) DEFAULT ''
);
ALTER TABLE transactions_accounts ADD CONSTRAINT transactions_accounts_accounttype_id_fkey
    FOREIGN KEY (accounttype_id) REFERENCES transactions_accounttype(accounttype_id);

CREATE INDEX transactions_accounts_accounttype_id_idx ON transactions_accounts (accounttype_id);

