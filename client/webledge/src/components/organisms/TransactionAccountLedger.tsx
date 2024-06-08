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
        let accountSign = String(formEntries.accountSign)
        let otherAccountSign = accountSign == "DEBIT"  ? "CREDIT" : "DEBIT";
        const otherAccountID = Number(formEntries.otherAccountID)
        let amount = Number(formEntries.amount)
        let dStr = String(formEntries.transactionDate)
        let txnDate: Date = new Date(dStr);
        // if amount is negative, swap debgits and credits
        console.log(amount, accountSign, otherAccountSign, txnDate)
        if (amount < 0 ) {
            amount = -amount
            accountSign = accountSign == "DEBIT"  ? "CREDIT" : "DEBIT";
            otherAccountSign = otherAccountSign == "DEBIT"  ? "CREDIT" : "DEBIT";
        }
        console.log(amount, accountSign, otherAccountSign, txnDate)
        // make debitAndCreditSet of two from this account and the selected account
        let dcSet: Array<TransactionDebitCreditRequest> = [
            {accountID: accountID, transactionDCAmount: amount, debitOrCredit: accountSign },
            {accountID: otherAccountID, transactionDCAmount: amount, debitOrCredit: otherAccountSign },
        ];
        const comment = String(formEntries.transactionComment)
        if (comment == ""){
            console.warn("comment cannot be empty")
        }
        const newTransaction : TransactionPostRequest = {
            transactionDate: txnDate.toISOString(),
            transactionComment: comment,
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
    let rowColor = "bg-slate-200"
    return (
        <div className="flex w-full flex-col md:col-span-4 grow justify-between rounded-xl bg-slate-100 p-4">
            <div className="text-xl font-bold">
                Add Transaction to {data?.accountFullName}
            </div>
            <div className="flex">
                <form className="flex" onSubmit={handleSubmit}>
                    <label className="my-4 text-xl font-bold mx-4 bg-slate-200">Date:
                        <input className="bg-slate-300 font-normal" type="date" name="transactionDate"/>
                    </label>
                    <label className="my-4 text-xl font-bold mx-4 bg-slate-200">Comment:
                        <input className="bg-slate-300 font-normal" type="text" name="transactionComment"/>
                    </label>
                    <label className="my-4 text-xl font-bold mx-4 bg-slate-200">Amount:
                        <input className="bg-slate-300" type="number" name="amount"/>
                    </label>
                    <label className="my-4 text-xl font-bold mx-4 bg-slate-200">
                        To Account:
                        <AccountSelector name={"otherAccountID"} id={0} includeTop={false}
                                         excludeID={data?.accountID}/>
                    </label>
                    <div className="bg-slate-300 flex">
                        <input className=" bg-slate-300" type="hidden" name="accountID"
                               defaultValue={data?.accountID}/>
                        <input className=" bg-slate-300" type="hidden" name="accountSign"
                               defaultValue={data?.accountSign}/>
                        <button className="p-3 font-bold" type="submit">Add Transaction</button>
                    </div>
                </form>
            </div>
            <div className="text-xl font-bold">
                Transactions
            </div>
            <div className="flex">
                <div className="w-8">
                    id
                </div>
                <div className="w-80">
                    Date
                </div>
                <div className="w-80">
                    Comment
                </div>
                <div className="w-80">
                    Account
                </div>
                <div className="w-16 text-right">
                    Amount
                </div>
                <div className="w-16 text-right">
                    Sign
                </div>
            </div>
            {data?.transactions && data.transactions.map((transaction: TransactionLedgerEntry, index: number) => {
            if (rowColor == "bg-slate-200"){
                rowColor = "bg-slate-300"
            } else {
                rowColor = "bg-slate-200"
            }
            let textColor = ""
            if (transaction.debitOrCredit != data.accountSign) {
                textColor = "text-red-500"
            }
            let txnDate: Date = new Date(transaction.transactionDate);
            return (
                <div className={'flex '+rowColor}  key={index}>
                    <Link to={{
                        pathname: '/transactions/' + transaction.transactionID,
                        search: '?returnAccount=' + accountID
                    }} className={`flex nav__item font-bold`}>
                        <div className="w-8">
                            {transaction.transactionID}
                        </div>
                        <div className="w-80">
                            {txnDate.toISOString().split('T')[0]}
                        </div>
                        <div className="w-80">
                            {transaction.transactionComment}
                        </div>
                        <div className="w-80">
                            {transaction.split}
                        </div>
                        <div className={"w-16 text-right "+textColor}>
                            {transaction.transactionDCAmount}
                        </div>
                        <div className={"w-16 text-right "+textColor}>
                            {transaction.debitOrCredit}
                        </div>
                    </Link>
                </div>
            );
        })}
        </div>
    );
}
