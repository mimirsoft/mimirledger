CREATE TYPE transaction_account_sign_type AS ENUM ('CREDIT','DEBIT');
CREATE TYPE transaction_account_type AS ENUM ('ASSET','LIABILITY','EQUITY','INCOME','EXPENSE','GAIN','LOSS');

CREATE TABLE transaction_accounts (
    account_id integer NOT NULL PRIMARY KEY,
    account_name varchar(50) DEFAULT '',
    account_memo varchar(70) DEFAULT '',
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
    account_code varchar(50) DEFAULT '',
    account_sign transaction_account_sign_type NOT NULL DEFAULT 'DEBIT',
    account_type transaction_account_type NOT NULL DEFAULT 'ASSET'
);


INSERT INTO transaction_accounts(account_id,
                                 account_name, account_memo,
                                 account_current, account_left, account_right,
                                 account_balance, account_subtotal,
                                 account_fullname,
                                 account_sign,
                                 account_type)
VALUES (1,
        'ASSETS','TOP LEVEL ASSETS SAMPLE',
        true, 1,2,
        0,0,
        'ASSETS', 'DEBIT', 'ASSET'),
       (2,
        'LIABILITIES','TOP LEVEL LIABILITY SAMPLE',
        true, 3,4,
        0,0,
        'ASSETS', 'CREDIT', 'LIABILITY'),
(3,
    'EQUITY','TOP LEVEL EQUITY SAMPLE',
    true, 5,6,
    0,0,
    'EQUITY', 'CREDIT', 'EQUITY')
