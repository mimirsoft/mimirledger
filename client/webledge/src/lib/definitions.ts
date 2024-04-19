export type TransactionAccountRequest = {
    accountParent: number;
    accountName: string;
    accountMemo: string;
    accountType: string;
};


export type TransactionAccount = {
    accountID: number;
    accountParent: number;
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
    transactionID: number;
};
