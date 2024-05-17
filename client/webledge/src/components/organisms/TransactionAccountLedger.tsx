import {useParams, useSearchParams, useNavigate, Link} from "react-router-dom";
import React, {FormEvent} from "react";
import {
    TransactionDebitCreditRequest, TransactionLedgerEntry,
    TransactionPostRequest
} from "../../lib/definitions";
import AccountSelector from "../molecules/AccountSelector";
import {useGetTransactionsOnAccountLedger} from "../../lib/data";

const postFormData = async (formData: FormData) => {
    try {
        // Do a bit of work to convert the entries to a plain JS object
        const myURL = new URL('/transactions', process.env.REACT_APP_MIMIRLEDGER_API_URL);

        // Do a bit of work to convert the entries to a plain JS object
        const formEntries = Object.fromEntries(formData);
        const accountID = Number(formEntries.accountID)
        const accountSign = String(formEntries.accountSign)
        const otherAccountSign = accountSign == "DEBIT"  ? "CREDIT" : "DEBIT";
        const otherAccountID = Number(formEntries.otherAccountID)
        const amount = Number(formEntries.amount)

        // make debitAndCreditSet of two from this account and the selected account
        let dcSet: Array<TransactionDebitCreditRequest> = [
            {accountID: accountID, transactionDCAmount: amount, debitOrCredit: accountSign },
            {accountID: otherAccountID, transactionDCAmount: amount, debitOrCredit: otherAccountSign },
        ];

        const newTransaction : TransactionPostRequest = {
            transactionComment : String(formEntries.transactionComment),
            debitCreditSet : dcSet,
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


export default function TransactionAccountLedger() {
    const { accountID } = useParams();
    const { data, isLoading, error } = useGetTransactionsOnAccountLedger(accountID);
    if (isLoading) return <div className="Loading">Loading...</div>
    if (error) return <div>Failed to load</div>

    return (
        <div className="flex w-full flex-col md:col-span-4">
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
                        <AccountSelector name={"otherAccountID"} id={0} includeTop={false} excludeID={data?.accountID}/>
                    </label>
                    <div className="bg-slate-300 flex">
                        <input className=" bg-slate-300" type="hidden" name="accountID" defaultValue={data?.accountID}/>
                        <input className=" bg-slate-300" type="hidden" name="accountSign"
                               defaultValue={data?.accountSign}/>
                        <button className="p-3 font-bold" type="submit">Add Transaction</button>
                    </div>
                </form>
            </div>
            <div className="flex">
                <div className="w-8">
                    id
                </div>
                <div className="w-80">
                    Comment
                </div>
                <div className="w-80">
                    Account
                </div>
                <div className="w-16">
                    Amount
                </div>
            </div>
            {data?.transactions && data.transactions.map((transaction: TransactionLedgerEntry, index: number) => {
                return (
                    <div className="flex" key={index}>
                        <div className="w-8">
                            {transaction.transactionID}
                        </div>
                        <div className="w-80">
                            {transaction.transactionComment}
                        </div>
                        <div className="w-80">
                            {transaction.split}
                        </div>
                        <div className="w-16 text-right">
                            {transaction.transactionDCAmount}
                        </div>
                        <Link to={{pathname: '/transactions/' + transaction.transactionID, search: '?returnAccount='+accountID}} className={`nav__item p-4 }`}>
                            EDIT Transaction
                        </Link>
                    </div>
                );
            })}
        </div>
    );
}
