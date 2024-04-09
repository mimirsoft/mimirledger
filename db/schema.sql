CREATE TYPE transaction_account_sign_type AS ENUM ('CREDIT','DEBIT');
CREATE TYPE transaction_account_type AS ENUM ('ASSET','LIABILITY','EQUITY','INCOME','EXPENSE','GAIN','LOSS');

CREATE TABLE transaction_accounts (
    account_id SERIAL PRIMARY KEY,
    account_parent integer NOT NULL DEFAULT '0',
    account_name varchar(50) DEFAULT '',
    account_full_name varchar(200) DEFAULT '',
    account_memo varchar(70) DEFAULT '',
    account_current bool NOT NULL DEFAULT true,
    account_left integer NOT NULL,
    account_right integer NOT NULL,
    account_balance integer NOT NULL DEFAULT 0,
    account_subtotal integer NOT NULL  DEFAULT 0,
    account_decimals smallint NOT NULL  DEFAULT 2,
    account_reconcile_date timestamp without time zone DEFAULT NULL,
    account_flagged bool NOT NULL DEFAULT false,
    account_locked bool NOT NULL DEFAULT false,
    account_open_date timestamp without time zone DEFAULT now(),
    account_close_date timestamp without time zone DEFAULT NULL,
    account_code varchar(50) DEFAULT NULL,
    account_sign transaction_account_sign_type NOT NULL DEFAULT 'DEBIT',
    account_type transaction_account_type NOT NULL DEFAULT 'ASSET'
);
CREATE INDEX transaction_accounts_account_left_idx ON transaction_accounts (account_left);
CREATE INDEX transaction_accounts_account_right_idx ON transaction_accounts (account_right);


CREATE TABLE transaction_main (
    transaction_id SERIAL PRIMARY KEY,
    transaction_date DATE NOT NULL DEFAULT CURRENT_DATE,
    transaction_comment varchar(250) NOT NULL CHECK (transaction_comment <> ''),
    transaction_amount integer NOT NULL CHECK (transaction_amount > 0),
    transaction_reference varchar(32) DEFAULT NULL,
    is_reconciled bool NOT NULL default FALSE,
    transaction_reconcile_date date DEFAULT NULL,
    is_split bool NOT NULL default FALSE) ;

CREATE TABLE transaction_debit_credit (
    transaction_dc_id SERIAL PRIMARY KEY,
    account_id integer NOT NULL,
    transaction_id integer NOT NULL,
    transaction_dc_amount decimal(11,2) DEFAULT '0.00',
    transaction_dc transaction_account_sign_type NOT NULL DEFAULT 'DEBIT') ;

ALTER TABLE transaction_debit_credit
    ADD CONSTRAINT transactions_debit_credit_transaction_id_fkey FOREIGN KEY (transaction_id) REFERENCES transaction_main(transaction_id) ON DELETE CASCADE;
CREATE INDEX transactions_debit_credit_transaction_id_idx ON transaction_debit_credit (transaction_id);
ALTER TABLE transaction_debit_credit
    ADD CONSTRAINT transactions_debit_credit_transaction_account_id_fkey FOREIGN KEY (account_id) REFERENCES transaction_accounts(account_id) ON DELETE CASCADE;
CREATE INDEX transactions_debit_credit_transaction_account_idx ON transaction_debit_credit (account_id);
