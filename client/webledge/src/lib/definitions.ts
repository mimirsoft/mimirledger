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
    accountMemo: string;
    accountBalance: number;
};

export type TransactionAccountType = {
    name:string
};