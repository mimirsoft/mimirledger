import { ArrowPathIcon } from '@heroicons/react/24/outline';
import clsx from 'clsx';
import { TransactionAccount, TransactionAccountRequest } from '../../lib/definitions';

import { useTransactionAccounts } from '../../lib/data';
import useSWR from "swr";
import React, {FormEvent} from "react";
const fetcher = (...args: Parameters<typeof fetch>) => fetch(...args).then(res => res.json())
const myURL = new URL('/accounts', process.env.REACT_APP_MIMIRLEDGER_API_URL);

const postFormData = async (formData: FormData) => {
    try {
        const myURL = new URL('/accounts', process.env.REACT_APP_MIMIRLEDGER_API_URL);
        console.log(myURL)
        // Do a bit of work to convert the entries to a plain JS object
        const formEntries = Object.fromEntries(formData);

        const newAccount : TransactionAccountRequest = {
            accountParent : Number(formEntries.accountParent),
            accountName : String(formEntries.accountName),
            accountType : String(formEntries.accountType),
            accountMemo : String(formEntries.accountMemo),
        };

        var json = JSON.stringify(newAccount);

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
                <select name="accountType">
                    <option value="ASSET">ASSET</option>
                    <option value="LIABILITY">LIABILITY</option>
                    <option value="EQUITY">EQUITY</option>
                    <option value="INCOME">INCOME</option>
                    <option value="ASSET">EXPENSE</option>
                    <option value="ASSET">GAIN</option>
                    <option value="ASSET">LOSS</option>
                </select>
            </label>
            <label>Memo:
                <input type="text" name="accountMemo"/>
            </label>

            <button type="submit">Create Account</button>
        </form>
    </div>
    <div>
        <div className="flex">
            <div className="w-80">
                AccountID
            </div>
            <div className="w-80">
                FullName
            </div>
            <div className="w-80">
                Name
            </div>
            <div className="w-80">
                Balance
            </div>
        </div>
        {data.accounts && data.accounts.map((account: TransactionAccount, index: number) => {
            return (
                <div className="flex">
                    <div className="w-80">
                        {account.accountID}
                    </div>
                    <div className="w-80">
                        {account.accountFullName}
                    </div>
                    <div className="w-80">
                        {account.accountName}
                    </div>
                    <div className="w-80">
                        {account.accountBalance}
                    </div>
                </div>
            );
        })}
    </div>
</div>

        </div>
    );
}
