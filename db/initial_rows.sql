
INSERT INTO transaction_accounts(account_id,
                                 account_name, account_memo,
                                 account_current, account_left, account_right,
                                 account_balance, account_subtotal,
                                 account_full_name,
                                 account_sign,
                                 account_type)
VALUES (1,
        'ASSETS','TOP LEVEL ASSETS SAMPLE',
        true, 1,2,
        0,0,
        'ASSETS FULLNAME', 'DEBIT', 'ASSET'),
       (2,
        'LIABILITIES','TOP LEVEL LIABILITY SAMPLE',
        true, 3,4,
        0,0,
        'LIABILITIES FULLNAME', 'CREDIT', 'LIABILITY'),
       (3,
        'EQUITY','TOP LEVEL EQUITY SAMPLE',
        true, 5,6,
        0,0,
        'EQUITY FULLNAME', 'CREDIT', 'EQUITY');

