import { useParams, useNavigate } from "react-router-dom";
import useSWR, { Fetcher } from "swr";
import React, {FormEvent} from "react";
import {TransactionAccount} from "../../lib/definitions";

const fetcher = (...args: Parameters<typeof fetch>) => fetch(...args).then(res => res.json())
const myURL = new URL('/transactions/account/', process.env.REACT_APP_MIMIRLEDGER_API_URL);

// get the transactions
const useTransactionOnAccount = (accountID:string |undefined):{data:TransactionAccount | undefined, isLoading:boolean, error: string|undefined} => {
    const { data, error , isLoading} = useSWR<TransactionAccount, string>(process.env.REACT_APP_MIMIRLEDGER_API_URL+'/transactions/account/'+accountID, fetcher);
    return {
        data,
        isLoading,
        error
    };
};
export default function TransactionLedger() {
    const { accountID } = useParams();
    const { data, isLoading, error } = useTransactionOnAccount(accountID);

    return (  <div className="flex w-full flex-col md:col-span-4">
        <h2 className={` mb-4 text-xl md:text-2xl`}>
            Transaction - Account
        </h2>
        </div>
        );
    }
