import { formatCurrency } from './utils';
import useSWR from "swr";
import {
    Account,
    AccountSet, TransactionAccountType, TransactionAccountTypeSet,
    TransactionLedgerResponse, TransactionResponse
} from "./definitions";


const accountURL = new URL('/accounts', process.env.REACT_APP_MIMIRLEDGER_API_URL);
export const  useGetAccounts = ():{data:AccountSet | undefined, isLoading:boolean, error: string|undefined} => {
   // const fetcher = (...args: Parameters<typeof fetch>) => fetch(...args).then(res => res.json())
    //console.log(fetcher);
   return useSWR<AccountSet, string>(accountURL);
    /*const { data, error, isLoading } = useSWR<AccountSet, string>(myURL);
    console.log("data"+data)
    return {
        data,
        isLoading,
        error,
    };
    */
}

export const useGetAccount = (accountID:string |undefined):{data:Account | undefined, isLoading:boolean, error: string|undefined} => {
    return useSWR<Account, string>(process.env.REACT_APP_MIMIRLEDGER_API_URL+'/accounts/'+accountID);
}

export const useGetTransaction = (transactionID:string |undefined):{data:TransactionResponse | undefined, isLoading:boolean, error: string|undefined} => {
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

// getAccounts - make map ID to name