
INSERT INTO transaction_accounts(account_name, account_memo,
                                 account_current, account_left, account_right,
                                 account_balance, account_subtotal,
                                 account_full_name,
                                 account_sign,
                                 account_type)
VALUES ('ASSETS','TOP LEVEL ASSETS SAMPLE',
        true, 1,2,
        0,0,
        'ASSETS', 'DEBIT', 'ASSET'),
       ('LIABILITIES','TOP LEVEL LIABILITY SAMPLE',
        true, 3,4,
        0,0,
        'LIABILITIES', 'CREDIT', 'LIABILITY'),
       ('EQUITY','TOP LEVEL EQUITY SAMPLE',
        true, 5,6,
        0,0,
        'EQUITY', 'CREDIT', 'EQUITY');

