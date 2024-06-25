import React, {FormEvent} from 'react'
import {Link, useParams, useSearchParams} from "react-router-dom";
import {useGetAccounts, useGetUnreconciledTransactionOnAccount} from "../../lib/data";
import {
    Account,
    TransactionLedgerEntry,
    TransactionReconciledPostRequest
} from "../../lib/definitions";
import {formatCurrency, formatCurrencyNoSign} from "../../lib/utils";

async function updateReconcileSearchDate(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const formData = new FormData(event.currentTarget)
    const formEntries = Object.fromEntries(formData);
    const searchDateString = String(formEntries.accountReconcileSearchDate)
    const accountID = Number(formEntries.accountID)
    window.open("/reconcile/"+accountID+"?date="+searchDateString,"_self");
};
const updateReconcileTransaction = async (formData: FormData) => {
    try {
        // Do a bit of work to convert the entries to a plain JS object
        const formEntries = Object.fromEntries(formData);
        const transactionID = Number(formEntries.transactionID)
        let dStr = String(formEntries.reconcileDate)
        let txnDate: Date = new Date(dStr);

        const myURL = new URL('/transactions/'+transactionID+"/reconciled", process.env.REACT_APP_MIMIRLEDGER_API_URL);

        const reconciledPostRequest : TransactionReconciledPostRequest = {
            transactionID : transactionID,
            transactionReconcileDate: txnDate.toISOString(),
        };
        var json = JSON.stringify(reconciledPostRequest);
        console.log(json);

        const settings :RequestInit = {
            method: 'PUT',
            body: json,
        };
        return await fetch(myURL, settings);
    }catch (error) {
        console.error('Error making POST request:', error);
    }
}
export default function AccountReconcileForm() {
    const {accountID} = useParams();
    let [searchParams] = useSearchParams();
    let date = searchParams.get("date");
    let searchDate = String(date)
    const { data:acctData, error:acctError, isLoading:acctIsLoading } = useGetAccounts()
    const {
        data,
        isLoading,
        error
    } = useGetUnreconciledTransactionOnAccount(accountID, searchDate);

    async function toggleReconciled(event: FormEvent<HTMLFormElement>) {
        event.preventDefault()
        const formData = new FormData(event.currentTarget)
        const response = await updateReconcileTransaction(formData);
    };

    if (isLoading) return <div className="Loading">Loading...</div>
    if (acctIsLoading) return <div className="Loading">Loading...</div>
    if (error) return <div>Failed to load</div>
    if (acctError) return <div>Failed to load</div>
    let acctMap = new Map<number, string>
    acctData?.accounts.map((acct: Account, index: number) => {
        acctMap.set(acct.accountID, acct.accountFullName)
    })

    let rowColor = "bg-slate-200"
    var minDate = new Date('0001-01-01T00:00:00Z');
    minDate.setDate(minDate.getDate() + 1);

    return (
        <div className="flex w-full flex-col md:col-span-4 grow justify-between rounded-xl bg-slate-100 p-4">
            <div className="flex font-bold text-xl">AccountReconcileForm</div>
            <div className="flex m-2">
                <form className="flex" onSubmit={updateReconcileSearchDate}>
                    <label className="my-4 text-xl font-bold mx-4 bg-slate-200">Date:
                        <input className="bg-slate-300 font-normal" type="date" name="accountReconcileSearchDate"
                               defaultValue={searchDate}/>
                    </label>
                    <div className="bg-slate-300 flex">
                        <input className=" bg-slate-300" type="hidden" name="accountID"
                               defaultValue={accountID}/>
                        <button className="p-3 font-bold" type="submit">Search For Date</button>
                    </div>
                </form>
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
                <div className="w-20 text-right mr-2">
                    Amount
                </div>
                <div className="w-16">
                    Sign
                </div>
                <div className="w-20">
                    Rec?
                </div>
                <div className="w-80">
                    ReconciledDate
                </div>
            </div>
            {data?.transactions && data.transactions.map((transaction: TransactionLedgerEntry, index: number) => {
                if (rowColor == "bg-slate-200") {
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
                let txnDate: Date = new Date(transaction.transactionDate);
                let txnReconciledDate: Date = new Date(transaction.transactionReconcileDate);

                let txnReconciledDateStr: string
                if (txnReconciledDate < minDate) {
                    txnReconciledDateStr = "";
                } else {
                    txnReconciledDateStr = txnReconciledDate.toISOString().split('T')[0]
                }
                let otherAccounts = [];
                let otherAccountStr = ""
                // if the transaction split has a comma, we have a split transaction
                if (transaction.split.indexOf(',') != -1) {
                    var segments = transaction.split.split(',');
                    for (let i = 0; i < segments.length; i++) {
                        // add to the array
                        otherAccounts.push(acctMap.get(Number(segments[i])))
                    }
                    otherAccountStr = otherAccounts.join(",")
                } else {
                    otherAccountStr = String(acctMap.get(Number(transaction.split)))
                }
                let txnReconciled = transaction.isReconciled ? "Y" : "N";
                return (
                    <div className={'flex ' + rowColor} key={index}>
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
                            {otherAccountStr}
                        </div>
                        <div className={"w-20 text-right mr-2 " + textColor}>
                            {formatCurrency(txnAmount)}
                        </div>
                        <div className={"w-16 " + textColor}>
                            {transaction.debitOrCredit}
                        </div>
                        <div className="w-20">
                            {txnReconciled}
                        </div>
                        <div className="w-80">
                            <form onSubmit={toggleReconciled}>
                                <input className="bg-slate-300 text-xl font-normal" type="date" name="reconcileDate"
                                       defaultValue={txnReconciledDateStr}/>
                                <input className=" bg-slate-300" type="hidden" name="transactionID"
                                       defaultValue={transaction.transactionID}/>
                                <button className="bg-blue-500 m-1 p-2 font-bold"
                                        type="submit">Mark Reconciled
                                </button>
                            </form>
                        </div>
                    </div>
                );
            })}
            <div className="flex m-2 justify-end">
                <label className="my-4 text-xl font-bold mx-4 bg-slate-200">Ending Balance:
                   <input className="w-24 text-xl bg-slate-300 text-right" type="text"
                       name="endingBalance"/>
                </label>

            </div>
        </div>
    );
}