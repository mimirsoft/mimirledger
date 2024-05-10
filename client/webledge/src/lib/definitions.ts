export type TransactionAccountPostRequest = {
    accountParent: number;
    accountName: string;
    accountMemo: string;
    accountType: string;
};

export type TransactionAccount = {
    accountID: number;
    accountParent: number;
    accountLeft: number;
    accountRight: number;
    accountName: string;
    accountFullName: string;
    accountType: string;
    accountMemo: string;
    accountBalance: number;
};

export type TransactionAccountType = {
    name:string
    sign:string
};


export type Transaction = {
    transactionID: number
    transactionComment: string
    transactionAmount: number
};

export type TransactionLedgerResponse = {
    accountID: number;
    accountName: string;
    accountFullName: string;
    accountSign: string;
    transactions: Transaction[];
};


export type TransactionPostRequest = {
    transactionComment: string
    debitCreditSet: string
};

