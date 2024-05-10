import { useParams, useNavigate } from "react-router-dom";
import useSWR, { Fetcher } from "swr";
import React, {FormEvent} from "react";
import {TransactionAccount, TransactionLedgerResponse, TransactionPostRequest} from "../../lib/definitions";
import AccountSelector from "../molecules/account-selector";
import AccountTypeSelector from "../molecules/account-type-selector";

const fetcher = (...args: Parameters<typeof fetch>) => fetch(...args).then(res => res.json())

const postFormData = async (formData: FormData) => {
    try {
        // Do a bit of work to convert the entries to a plain JS object
        const myURL = new URL('/transactions', process.env.REACT_APP_MIMIRLEDGER_API_URL);
        console.log(myURL)

        // Do a bit of work to convert the entries to a plain JS object
        const formEntries = Object.fromEntries(formData);
        const accountID = Number(formEntries.accountID)
        console.log(formEntries.accountParent)
        // make debitAndCreditSet of two

        const newTransaction : TransactionPostRequest = {
            transactionComment : String(formEntries.transactionComment),
            debitCreditSet :"",
        };
        var json = JSON.stringify(newTransaction);
        console.log(json)
        const settings :RequestInit = {
            method: 'POST',
            body: json,
        };
        return await fetch(myURL, settings);
    } catch (error) {
        console.error('Error making POST request:', error);
    }
}
async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const formData = new FormData(event.currentTarget)
    const result = await postFormData(formData);
    window.location.reload();
};
// get the transactions
const useTransactionOnAccount = (accountID:string |undefined):{data:TransactionLedgerResponse | undefined, isLoading:boolean, error: string|undefined} => {
    const { data,
        error ,
        isLoading} = useSWR<TransactionLedgerResponse, string>(process.env.REACT_APP_MIMIRLEDGER_API_URL+'/transactions/account/'+accountID, fetcher);
    return {
        data,
        isLoading,
        error
    };
};
export default function TransactionAccountLedger() {
    const { accountID } = useParams();
    const { data, isLoading, error } = useTransactionOnAccount(accountID);
    if (isLoading) return <div className="Loading">Loading...</div>
    if (error) return <div>Failed to load</div>

    console.log(data)
    return (<div className="flex w-full flex-col md:col-span-4">
            <h2 className={` mb-4 text-xl md:text-2xl`}>
                Add Transaction to {data?.accountFullName}
            </h2>
            <div className="flex">
                <form className="flex" onSubmit={handleSubmit}>
                    <label className="my-4 text-xl font-bold mx-4 bg-slate-200">Comment:
                        <input className="bg-slate-300 font-normal" type="text" name="transactionComment"/>
                    </label>
                    <label className="my-4 text-xl font-bold mx-4 bg-slate-200">Amount:
                        <input className="bg-slate-300" type="text" name="amount"/>
                    </label>
                    <label className="my-4 text-xl font-bold mx-4 bg-slate-200">
                        To Account:
                        <AccountSelector id={0} includeTop={false} excludeID={data?.accountID}/>
                    </label>
                    <div className="bg-slate-300 flex">
                        <input className=" bg-slate-300" type="hidden" name="accountID" defaultValue={data?.accountID}/>
                        <input className=" bg-slate-300" type="hidden" name="accountSign"
                               defaultValue={data?.accountSign}/>
                        <button className="p-3 font-bold" type="submit">Add Transaction</button>
                    </div>
                </form>
            </div>

        </div>
    );
}
