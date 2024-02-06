CREATE TYPE transaction_account_sign_type AS ENUM ('CREDIT','DEBIT');
CREATE TYPE transaction_account_type AS ENUM ('ASSET','LIABILITY','EQUITY','INCOME','EXPENSE','GAIN','LOSS');

CREATE TABLE transaction_accounts (
    account_id SERIAL PRIMARY KEY,
    account_parent integer NOT NULL DEFAULT '0',
    account_name varchar(50) DEFAULT '',
    account_full_name varchar(200) DEFAULT '',
    account_memo varchar(70) DEFAULT '',
    account_current bool NOT NULL DEFAULT true,
    account_left integer DEFAULT NULL,
    account_right integer DEFAULT NULL,
    account_balance decimal(12,2) DEFAULT NULL,
    account_subtotal decimal(12,2) DEFAULT NULL,
    account_reconcile_date timestamp without time zone DEFAULT NULL,
    account_flagged bool NOT NULL DEFAULT false,
    account_locked bool NOT NULL DEFAULT false,
    account_open_date timestamp without time zone DEFAULT now(),
    account_close_date timestamp without time zone DEFAULT NULL,
    account_code varchar(50) DEFAULT NULL,
    account_sign transaction_account_sign_type NOT NULL DEFAULT 'DEBIT',
    account_type transaction_account_type NOT NULL DEFAULT 'ASSET'
);


CREATE TABLE transactions_main (
    transaction_id integer NOT NULL PRIMARY KEY,
    transaction_date DATE NOT NULL DEFAULT CURRENT_DATE,
    transaction_comment varchar(250) DEFAULT NULL,
    transaction_amount decimal(11,2) DEFAULT '0.00',
    transaction_check_num varchar(32) DEFAULT NULL,
    transaction_reconcile bool NOT NULL default FALSE,
    transaction_reconcile_date date DEFAULT NULL,
    is_split bool NOT NULL default FALSE) ;

CREATE TABLE transactions_debit_credit (
    transaction_dc_id integer NOT NULL PRIMARY KEY,
    account_id integer NOT NULL,
    transaction_id integer NOT NULL,
    transaction_dc_amount decimal(11,2) DEFAULT '0.00',
    transaction_dc transaction_account_sign_type NOT NULL DEFAULT 'DEBIT') ;

ALTER TABLE transactions_debit_credit
    ADD CONSTRAINT transactions_debit_credit_transaction_id_fkey FOREIGN KEY (transaction_id) REFERENCES transactions_main(transaction_id) ON DELETE CASCADE;
CREATE INDEX transactions_debit_credit_transaction_id_idx ON transactions_debit_credit (transaction_id);
ALTER TABLE transactions_debit_credit
    ADD CONSTRAINT transactions_debit_credit_transaction_account_id_fkey FOREIGN KEY (account_id) REFERENCES transaction_accounts(account_id) ON DELETE CASCADE;
CREATE INDEX transactions_debit_credit_transaction_account_idx ON transactions_debit_credit (account_id);
