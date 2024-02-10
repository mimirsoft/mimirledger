import { formatCurrency } from './utils';
import useSWR from "swr";

export function fetchTransactionAccounts():{data: any, error: any, isValidating:boolean} {
    const fetcher = (...args) => fetch(...args).then((res) => res.json())

    console.log(process.env.NEXT_PUBLIC_MIMIRLEDGER_API_URL);
    const myURL = new URL('/accounts', process.env.NEXT_PUBLIC_MIMIRLEDGER_API_URL);
    console.log(myURL);
   // const { data, error, isValidating } = useSWR(myURL, fetcher)
    return useSWR(myURL, fetcher)
    //return data, error, isValidating;
}
