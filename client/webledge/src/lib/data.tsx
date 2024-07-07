import useSWR from "swr";
import {
    Account, AccountReconcileResponse,
    AccountSet, TransactionAccountType, TransactionAccountTypeSet,
    TransactionLedgerResponse, TransactionResponse,
    ReportSet, Report
} from "./definitions";
import {KeyedMutator} from "swr/_internal";


export const  useGetAccounts = ():{data:AccountSet | undefined, isLoading:boolean, error: string|undefined} => {
   return useSWR<AccountSet, string>(process.env.REACT_APP_MIMIRLEDGER_API_URL+'/accounts');
}

export const useGetAccount = (accountID:string |undefined):{data:Account | undefined, isLoading:boolean, error: string|undefined} => {
    return useSWR<Account, string>(process.env.REACT_APP_MIMIRLEDGER_API_URL+'/accounts/'+accountID);
}

export const useGetTransaction = (transactionID:string |undefined):{
    data:TransactionResponse | undefined,
    isLoading:boolean, error: string|undefined
    mutate: KeyedMutator<TransactionResponse> } => {
    return useSWR<TransactionResponse, string>(process.env.REACT_APP_MIMIRLEDGER_API_URL+'/transactions/'+transactionID);
}

// get the transactionLedger
export const useGetTransactionsOnAccountLedger = (accountID:string |undefined):{data:TransactionLedgerResponse | undefined, isLoading:boolean, error: string|undefined} => {
    return useSWR<TransactionLedgerResponse, string>(process.env.REACT_APP_MIMIRLEDGER_API_URL+'/transactions/account/'+accountID);
};

const accountTypesURL = new URL('/accounttypes', process.env.REACT_APP_MIMIRLEDGER_API_URL);
export const useGetTransactionAccountTypes =  ():{data:TransactionAccountTypeSet | undefined, isLoading:boolean, error: string|undefined} => {
   return useSWR(accountTypesURL)
}

export const  useGetReports = ():{data:ReportSet | undefined, isLoading:boolean, error: string|undefined} => {
    return useSWR<ReportSet, string>(process.env.REACT_APP_MIMIRLEDGER_API_URL+'/reports');
}
export const useGetReport = (reportID:string |undefined):{data:Report | undefined, isLoading:boolean, error: string|undefined} => {
    return useSWR<Report, string>(process.env.REACT_APP_MIMIRLEDGER_API_URL+'/reports/'+reportID);
}

// get the transactionLedger
export const useGetUnreconciledTransactionOnAccount = (accountID:string |undefined, date:string):{
    data:AccountReconcileResponse | undefined,
    isLoading:boolean,
    mutate: KeyedMutator<AccountReconcileResponse>,
    error: string|undefined} => {
    return useSWR<AccountReconcileResponse, string>(process.env.REACT_APP_MIMIRLEDGER_API_URL+
        '/transactions/account/'+accountID+'/unreconciled?date='+date);
};
