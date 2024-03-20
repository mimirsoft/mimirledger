import { ArrowPathIcon } from '@heroicons/react/24/outline';
import clsx from 'clsx';
import { TransactionAccount } from '../../lib/definitions';

import { useTransactionAccounts } from '../../lib/data';
import useSWR from "swr";
import React, {FormEvent} from "react";
const fetcher = (...args: Parameters<typeof fetch>) => fetch(...args).then(res => res.json())
const myURL = new URL('/accounts', process.env.REACT_APP_MIMIRLEDGER_API_URL);

const postFormData = async (formData: FormData) => {
    try {
        const myURL = new URL('/accounts', process.env.REACT_APP_MIMIRLEDGER_API_URL);
        console.log(myURL)
        var json =JSON.stringify(Object.fromEntries(formData));
        console.log(json)
        const settings :RequestInit = {
            method: 'POST',
            body: json,
        };
        const response = await fetch(myURL, settings);
        const result = await response.json();
        console.log('POST request result:', result);
    } catch (error) {
        console.error('Error making POST request:', error);
    }
}
async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const formData = new FormData(event.currentTarget)
    const result = await postFormData(formData);
};
export default function TransactionAccounts(){

    const { data, error, isLoading } = useSWR(myURL, fetcher)

    if (isLoading) return <div className="Loading">Loading...</div>
    if (error) return <div>Failed to load</div>

    return (
        <div className="flex w-full flex-col md:col-span-4">
<h2 className={` mb-4 text-xl md:text-2xl`}>
My Accounts
</h2>
<div className="flex grow flex-col justify-between rounded-xl bg-gray-50 p-4">
    <div className="bg-white px-6">
        <form onSubmit={handleSubmit}>
            <label>AccountName:
                <input type="text" name="accountName"/>
            </label>
            <label>
                AccountParent:
                <select name="accountParent">
                    <option value="0">Top Level</option>
                    {data.accounts && data.accounts.map((account: TransactionAccount, index: number) => {
                        return (
                            <option value="{account.accountID}"> {account.accountFullName}</option>
                        );
                    })}
                </select>
            </label>
            <label>AccountType:
                <input type="text" name="accountType"/>
            </label>
            <button type="submit">Create Account</button>
        </form>
        {data.accounts && data.accounts.map((account: TransactionAccount, index: number) => {
            return (
                <div>
                    <div className="flex items-center">
                        <div className="min-w-0">

                                <p className="hidden text-sm text-gray-500 sm:block">
                                    {account.accountFullName}
                                </p>
                            </div>
                        </div>
                        <p
                            className={` truncate text-sm font-medium md:text-base`}
                        >
                            {account.accountBalance}
                        </p>
                    </div>
                );
            })}
    </div>
</div>

        </div>
);
}
