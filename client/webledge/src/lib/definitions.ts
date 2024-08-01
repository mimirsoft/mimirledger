export type AccountSet = {
    accounts: Account[];
};

export type Account = {
    accountID: number;
    accountParent: number;
    accountLeft: number;
    accountRight: number;
    accountName: string;
    accountFullName: string;
    accountType: string;
    accountSign: string;
    accountMemo: string;
    accountBalance: number;
    accountReconcileDate: string;

};

export type TransactionAccountPostRequest = {
    accountParent: number;
    accountName: string;
    accountMemo: string;
    accountType: string;
};

export type TransactionAccountTypeSet = {
    accountTypes:TransactionAccountType[]
};

export type TransactionAccountType = {
    name:string
    sign:string
};

export type TransactionResponse = {
    transactionID: number
    transactionDate: string,
    transactionComment: string
    transactionAmount: number
    debitCreditSet: TransactionDebitCreditResponse[]
};

export type TransactionDebitCreditResponse = {
    transactionID: number
    accountID: number
    transactionDCAmount: number
    debitOrCredit: string
};

export type TransactionPostRequest = {
    transactionDate: string
    transactionComment: string
    debitCreditSet: TransactionDebitCreditRequest[]
};

export type TransactionEditPostRequest = {
    transactionID: number
    transactionDate: string
    transactionComment: string
    debitCreditSet: TransactionDebitCreditRequest[]
};

export type TransactionReconciledPostRequest = {
    transactionID: number
    transactionReconcileDate: string
};


export type TransactionDebitCreditRequest = {
    accountID: number
    transactionDCAmount: number
    debitOrCredit: string
};

export type TransactionLedgerResponse = {
    accountID: number;
    accountName: string;
    accountFullName: string;
    accountSign: string;
    transactions: TransactionLedgerEntry[];
};

export type TransactionLedgerEntry = {
    transactionID: number
    transactionDate: string,
    transactionReconcileDate: string
    transactionComment: string
    split: string
    transactionDCAmount: number
    debitOrCredit: string
    isSplit: boolean
    isReconciled: boolean
};
export type AccountReconcileDatePostRequest = {
    accountID: number
    accountReconcileDate: string
};

export type AccountReconcileResponse = {
    accountID: number;
    accountReconcileDate: string;
    priorReconciledBalance: number;
    searchDate: string;
    accountName: string;
    accountFullName: string;
    accountSign: string;
    transactions: TransactionLedgerEntry[];
};

export type ReportSet = {
    reports: Report[];
};

export type Report = {
    reportID: number;
    reportName: string;
    reportBody: ReportBody;
};
export type ReportBody = {
    sourceAccountSetType: string;
    sourceAccountGroup: string;
    sourcePredefinedAccounts: number[];
    sourceRecurseSubAccounts: boolean;
    sourceRecurseSubAccountsDepth: number;
    filterAccountSetType: string;
    filterAccountGroup: string;
    filterPredefinedAccounts: number[];
    filterRecurseSubAccounts: boolean;
    filterRecurseSubAccountsDepth: number;
    dataSetType: string;
};

export type ReportPostRequest = {
    reportName: string;
    reportBody: ReportBody;
};
