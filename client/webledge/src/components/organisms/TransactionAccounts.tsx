import {
    Account,
    TransactionAccountPostRequest,
} from '../../lib/definitions';
import {useState} from "react";
import Modal from '../molecules/Modal'
import React, {FormEvent} from "react";
import {Link, useSearchParams} from "react-router-dom";
import AccountSelector from "../molecules/AccountSelector";
import AccountTypeSelector from "../molecules/AccountTypeSelector";
import {useGetAccounts} from "../../lib/data";
import {formatCurrency} from "../../lib/utils";

const postFormData = async (formData: FormData) => {
    try {
        const myURL = new URL('/accounts', process.env.REACT_APP_MIMIRLEDGER_API_URL);
        console.log(myURL)
        // Do a bit of work to convert the entries to a plain JS object
        const formEntries = Object.fromEntries(formData);
        console.log(formEntries.accountParent)

        const newAccount : TransactionAccountPostRequest = {
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


export default function TransactionAccounts(){
    const { data, error, isLoading } = useGetAccounts()

    let [searchParams] = useSearchParams();
    let returnAccountID = searchParams.get("returnAccount");

    if (isLoading) return <div className="Loading">Loading...</div>
    if (error) return <div>Failed to load</div>
    let rowColor = "bg-slate-200"
    var minDate = new Date('0001-01-01T00:00:00Z');
    minDate.setDate(minDate.getDate() + 1);
    let reconcileDate: Date = new Date();
    let reconcileDateStr = reconcileDate.toISOString().split('T')[0]

    return (
        <div className="flex w-full flex-col md:col-span-4">
            <div className="flex grow flex-col justify-between rounded-xl bg-slate-100 p-4">
                <div className="text-xl font-bold">
                    Create New Account
                </div>
                <div className="flex">
                    <form className="flex" onSubmit={handleSubmit}>
                        <label className="my-4 mr-4 text-xl font-bold bg-slate-200">AccountName:
                            <input className="bg-slate-300 font-normal" type="text" name="accountName"/>
                        </label>
                        <label className="my-4 text-xl font-bold mr-4 bg-slate-200">
                            AccountParent:
                            <AccountSelector name={"accountParent"} id={0} includeTop={true} excludeID={0}
                                 multiple={false}
                                 multiSize={1}/>
                        </label>
                        <label className="my-4 mr-4 text-xl font-bold">AccountType:
                            <AccountTypeSelector selectedName=""/>
                        </label>
                        <label className="my-4 text-xl font-bold mr-4 bg-slate-200">Memo:
                            <input className=" bg-slate-300" type="text" name="accountMemo"/>
                        </label>
                        <div className="bg-slate-300 flex">
                            <button className="p-3 font-bold" type="submit">Create Account</button>
                        </div>
                    </form>
                </div>
                <div className="text-xl font-bold">
                    My Accounts
                </div>
                <div className="flex">
                    <div className="w-8 font-bold">
                        ID
                    </div>
                    <div className="w-8 font-bold">
                        L
                    </div>
                    <div className="w-8 font-bold">
                        R
                    </div>
                    <div className="w-80 font-bold">
                        FullName
                    </div>
                    <div className="w-32 font-bold text-right mr-4">
                        Balance
                    </div>
                    <div className="w-80 font-bold">
                        Name
                    </div>
                    <div className="w-20 font-bold">
                        Type
                    </div>
                    <div className="w-20 font-bold">
                        Sign
                    </div>
                    <div className="w-24 font-bold">
                        Rec Date
                    </div>
                </div>
                {data?.accounts && data.accounts.map((account: Account, index: number) => {
                    if (rowColor == "bg-slate-200"){
                        rowColor = "bg-slate-400"
                    } else {
                        rowColor = "bg-slate-200"
                    }
                    let textColor = ""
                    if (account.accountBalance < 0) {
                        textColor = "text-red-500"
                    }
                    let acctReconciledDate: Date = new Date(account.accountReconcileDate);

                    let acctReconciledDateStr :string
                    if (acctReconciledDate < minDate) {
                        acctReconciledDateStr = "";
                    } else {
                        acctReconciledDateStr = acctReconciledDate.toISOString().split('T')[0]
                    }
                    return (
                        <div className={'flex '+rowColor} key={index}>
                            <Link to={'/transactions/account/' + account.accountID}
                                  className={`flex mr-4 }`}>
                                <div className="w-8">
                                    {account.accountID}
                                </div>
                                <div className="w-8">
                                    {account.accountLeft}
                                </div>
                                <div className="w-8">
                                    {account.accountRight}
                                </div>
                                <div className="w-80">
                                    {account.accountFullName}
                                </div>
                                <div className={"w-32 text-right mr-4 " + textColor}>
                                    {formatCurrency(account.accountBalance)}
                                </div>
                                <div className="w-80">
                                    {account.accountName}
                                </div>
                                <div className="w-20">
                                    {account.accountType}
                                </div>
                                <div className="w-20">
                                    {account.accountSign}
                                </div>
                                <div className="w-24">
                                    {acctReconciledDateStr}
                                </div>
                            </Link>
                            <Link to={'/accounts/' + account.accountID} className={`nav__item font-bold mr-2`}>
                                Edit Account
                            </Link>
                            <Link to={{
                                pathname: '/reconcile/' + account.accountID,
                                search: '?date=' + reconcileDateStr
                                }} className={`font-bold`}>
                                Reconcile
                            </Link>
                        </div>
                    );
                })}
            </div>
        </div>
    );
}
