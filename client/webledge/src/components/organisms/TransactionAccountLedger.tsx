import {useParams, useSearchParams, useNavigate, Link} from "react-router-dom";
import React, {FormEvent} from "react";
import {
    Account,
    TransactionDebitCreditRequest, TransactionLedgerEntry,
    TransactionPostRequest
} from "../../lib/definitions";
import AccountSelector from "../molecules/AccountSelector";
import {useGetAccounts, useGetTransactionsOnAccountLedger} from "../../lib/data";
import {formatCurrency, parseCurrency} from "../../lib/utils";
const postFormData = async (formData: FormData) => {
    try {
        // Do a bit of work to convert the entries to a plain JS object
        const myURL = new URL('/transactions', import.meta.env.VITE_APP_MIMIRLEDGER_API_URL);

        // Do a bit of work to convert the entries to a plain JS object
        const formEntries = Object.fromEntries(formData);
        const accountID = Number(formEntries.accountID)
        let accountSign = String(formEntries.accountSign)
        let otherAccountSign = accountSign == "DEBIT"  ? "CREDIT" : "DEBIT";
        const otherAccountID = Number(formEntries.otherAccountID)
        let amount = parseCurrency(formEntries.amount)
        const dStr = String(formEntries.transactionDate)
        const txnDate: Date = new Date(dStr);
        // if amount is negative, swap debgits and credits
        console.log(amount, accountSign, otherAccountSign, txnDate)
        if (amount < 0 ) {
            amount = -amount
            accountSign = accountSign == "DEBIT"  ? "CREDIT" : "DEBIT";
            otherAccountSign = otherAccountSign == "DEBIT"  ? "CREDIT" : "DEBIT";
        }
        console.log(amount, accountSign, otherAccountSign, txnDate)
        // make debitAndCreditSet of two from this account and the selected account
        const dcSet: Array<TransactionDebitCreditRequest> = [
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
        const json = JSON.stringify(newTransaction);
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
}


export default function TransactionAccountLedger() {
    const { accountID } = useParams();

    const { data:acctData, error:acctError, isLoading:acctIsLoading } = useGetAccounts()
    const { data, isLoading, error } = useGetTransactionsOnAccountLedger(accountID);

    if (isLoading) return <div className="Loading">Loading...</div>
    if (acctIsLoading) return <div className="Loading">Loading...</div>
    if (error) return <div>Failed to load</div>
    if (acctError) return <div>Failed to load</div>

    const acctMap = new Map<number, string>
    acctData?.accounts.map((acct: Account, index: number) => {
        acctMap.set(acct.accountID, acct.accountFullName)
    })
    const minDate = new Date('0001-01-01T00:00:00Z');
    minDate.setDate(minDate.getDate() + 1);

    let runningTotal = 0
    console.log(acctMap)
    let rowColor = "bg-slate-200"

    const reconcileDate: Date = new Date();
    const reconcileDateStr = reconcileDate.toISOString().split('T')[0]

    return (
        <div className="flex w-full flex-col md:col-span-4 grow justify-between rounded-xl bg-slate-100 p-4">
            <div className="flex text-xl font-bold">
                <div className="mr-2">Add Transaction to {data?.accountFullName}</div>
                <Link to={{
                    pathname: '/reconcile/' + accountID,
                    search: '?date=' + reconcileDateStr
                }} className="underline text-blue-600 hover:text-blue-800 visited:text-purple-600">
                    Reconcile
                </Link>
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
                             excludeID={data?.accountID}
                             multiple={false}
                             multiSize={1}/>
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
                <div className="w-24">
                    Date
                </div>
                <div className="flex-grow">
                    Comment
                </div>
                <div className="w-80">
                    Account
                </div>
                <div className="w-20 text-right mr-2">
                    Amount
                </div>
                <div className="w-16">
                    Sign
                </div>
                <div className="w-8">
                    Rec
                </div>
                <div className="w-24">
                    Rec Date
                </div>
                <div className="w-20">
                    Balance
                </div>
            </div>
            {data?.transactions && data.transactions.map((transaction: TransactionLedgerEntry, index: number) => {
                if (rowColor == "bg-slate-200"){
                    rowColor = "bg-slate-300"
                } else {
                    rowColor = "bg-slate-200"
                }
                console.log(transaction)
                let textColor = ""
                let txnAmount = transaction.transactionDCAmount
                if (transaction.debitOrCredit != data.accountSign) {
                    textColor = "text-red-500"
                    txnAmount = -txnAmount
                }
                runningTotal +=txnAmount
                let runningTotalColor = ""
                if (runningTotal < 0) {
                    runningTotalColor = "text-red-500"
                }

                const txnDate: Date = new Date(transaction.transactionDate);
                const txnReconciledDate: Date = new Date(transaction.transactionReconcileDate);

                let txnReconciledDateStr :string
                if (txnReconciledDate < minDate) {
                    txnReconciledDateStr = "";
                } else {
                    txnReconciledDateStr = txnReconciledDate.toISOString().split('T')[0]
                }

                const otherAccounts= [];
                let otherAccountStr = ""
                // if the transaction split has a comma, we have a split transaction
                if (transaction.split.indexOf(',') != -1) {
                    const segments = transaction.split.split(',');
                    for(let i=0; i<segments.length; i++){
                        // add to the array
                        otherAccounts.push(acctMap.get(Number(segments[i])))
                    }
                    otherAccountStr = otherAccounts.join(",")
                } else {
                    otherAccountStr = String(acctMap.get(Number(transaction.split)))
                }
                const txnReconciled = transaction.isReconciled ? "Y" : "N";

                return (
                    <Link to={{
                        pathname: '/transactions/' + transaction.transactionID,
                        search: '?returnAccount=' + accountID
                    }} className={`flex w-full font-bold`}>
                        <div className={'flex flex-grow '+rowColor}  key={index}>
                            <div className="w-8">
                                {transaction.transactionID}
                            </div>
                            <div className="w-24">
                                {txnDate.toISOString().split('T')[0]}
                            </div>
                            <div className="flex-grow">
                                {transaction.transactionComment}
                            </div>
                            <div className="w-80">
                                {otherAccountStr}
                            </div>
                            <div className={"w-20 text-right mr-2 " + textColor}>
                                {formatCurrency(txnAmount)}
                            </div>
                            <div className={"w-16 " + textColor}>
                                {transaction.debitOrCredit}
                            </div>
                            <div className="w-8">
                                {txnReconciled}
                            </div>
                            <div className="w-24">
                                {txnReconciledDateStr}
                            </div>
                            <div className={"w-20 text-right font-bold " + runningTotalColor}>
                                {formatCurrency(runningTotal)}
                            </div>
                    </div>
                    </Link>

                );
            })}
        </div>
    );
}
