import useSWR from "swr";
import {
    Account, AccountReconcileResponse,
    AccountSet, TransactionAccountTypeSet,
    TransactionLedgerResponse, TransactionResponse,
    ReportSet, Report, ReportOutput
} from "./definitions";
import {KeyedMutator} from "swr/_internal";


export const  useGetAccounts = ():{data:AccountSet | undefined, isLoading:boolean, error: string|undefined} => {
   return useSWR<AccountSet, string>(import.meta.env.VITE_APP_SERVER_API_URL+'/accounts');
}

export const useGetAccount = (accountID:string |undefined):{data:Account | undefined, isLoading:boolean, error: string|undefined} => {
    return useSWR<Account, string>(import.meta.env.VITE_APP_SERVER_API_URL+'/accounts/'+accountID);
}
export const useGetReportOutput = (reportID:string |undefined,startDate:string,endDate:string):{
    data:ReportOutput | undefined, isLoading:boolean, error: string|undefined} => {
    return useSWR<ReportOutput, string>(import.meta.env.VITE_APP_SERVER_API_URL+'/reports/'+reportID+
        '/output?startDate='+
        startDate+'&endDate='+endDate);
}

export const useGetTransaction = (transactionID:string |undefined):{
    data:TransactionResponse | undefined,
    isLoading:boolean, error: string|undefined
    mutate: KeyedMutator<TransactionResponse> } => {
    return useSWR<TransactionResponse, string>(import.meta.env.VITE_APP_SERVER_API_URL+'/transactions/'+transactionID);
}

// get the transactionLedger
export const useGetTransactionsOnAccountLedger = (accountID:string |undefined):{data:TransactionLedgerResponse | undefined, isLoading:boolean, error: string|undefined} => {
    return useSWR<TransactionLedgerResponse, string>(import.meta.env.VITE_APP_SERVER_API_URL+'/transactions/account/'+accountID);
};

const accountTypesURL = new URL('/accounttypes', import.meta.env.VITE_APP_SERVER_API_URL);
export const useGetTransactionAccountTypes =  ():{data:TransactionAccountTypeSet | undefined, isLoading:boolean, error: string|undefined} => {
   return useSWR(accountTypesURL)
}

export const  useGetReports = ():{data:ReportSet | undefined, isLoading:boolean, error: string|undefined} => {
    return useSWR<ReportSet, string>(import.meta.env.VITE_APP_SERVER_API_URL+'/reports');
}
export const useGetReport = (reportID:string |undefined):{data:Report | undefined, isLoading:boolean, error: string|undefined} => {
    return useSWR<Report, string>(import.meta.env.VITE_APP_SERVER_API_URL+'/reports/'+reportID);
}

// get the transactionLedger
export const useGetUnreconciledTransactionOnAccount = (accountID:string |undefined, date:string):{
    data:AccountReconcileResponse | undefined,
    isLoading:boolean,
    mutate: KeyedMutator<AccountReconcileResponse>,
    error: string|undefined} => {
    return useSWR<AccountReconcileResponse, string>(import.meta.env.VITE_APP_SERVER_API_URL+
        '/transactions/account/'+accountID+'/unreconciled?date='+date);
};
