import { formatCurrency } from './utils';
import useSWR from "swr";
const fetcher = (...args: Parameters<typeof fetch>) => fetch(...args).then(res => res.json())

export function useTransactionAccounts() {

    const myURL = new URL('/accounts', process.env.REACT_APP_MIMIRLEDGER_API_URL);
    console.log(myURL);
    const { data, error, isLoading } = useSWR(myURL, fetcher)
    console.log(data);

    return {
        data: data,
        isLoading,
        isError: error
    }
}

export function useTransactionAccountTypes() {

    const myURL = new URL('/accounttypes', process.env.REACT_APP_MIMIRLEDGER_API_URL);
    console.log(myURL);
    const { data, error, isLoading } = useSWR(myURL, fetcher)
    console.log(data);

    return {
        data: data,
        isLoading,
        isError: error
    }
}
