import React, {FormEvent} from "react";
import {AccountReconcileResponse, TransactionReconciledPostRequest} from "../../lib/definitions";
import {KeyedMutator} from "swr/_internal";

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
};
const updateUnreconcileTransaction = async (formData: FormData) => {
    try {
        // Do a bit of work to convert the entries to a plain JS object
        const formEntries = Object.fromEntries(formData);
        const transactionID = Number(formEntries.transactionID)

        const myURL = new URL('/transactions/'+transactionID+"/unreconciled", process.env.REACT_APP_MIMIRLEDGER_API_URL);


        const settings :RequestInit = {
            method: 'PUT',
         };
        return await fetch(myURL, settings);
    }catch (error) {
        console.error('Error making POST request:', error);
    }
};

export default function TransactionToggleReconcileForm(props:{
    transactionID: number|undefined;
    reconciledDate: string;
    mutator:KeyedMutator<AccountReconcileResponse>,
    isReconciled:boolean}){
    async function toggleReconciled(event: FormEvent<HTMLFormElement>) {
        event.preventDefault()
        const formData = new FormData(event.currentTarget)
        const response = await updateReconcileTransaction(formData);
        if (response?.status == 200){
            await props.mutator()
        }
        else {
            console.log("ERROR"+response)
        }
    };

    async function toggleUnreconciled(event: FormEvent<HTMLFormElement>) {
        event.preventDefault()
        const formData = new FormData(event.currentTarget)
        const response = await updateUnreconcileTransaction(formData);
        if (response?.status == 200){
            await props.mutator()
        }
        else {
            console.log("ERROR"+response)
        }
    };

    if (props.isReconciled) {
        return (
            <form key={'toggleUnrec'+props.transactionID} onSubmit={toggleUnreconciled}>
                <input className="w-48 bg-slate-300 text-xl font-normal" type="date" name="reconcileDate"
                       defaultValue={props.reconciledDate}/>
                <input className="bg-slate-300" type="hidden" name="transactionID"
                       defaultValue={props.transactionID}/>
                <button className="bg-red-600 m-1 p-2 font-bold"
                        type="submit">Unreconcile
                </button>
            </form>
        )
    }
    return (
        <form key={'toggleRec'+props.transactionID} onSubmit={toggleReconciled}>
            <input className="w-48 bg-slate-300 text-xl font-normal" type="date" name="reconcileDate"
                   defaultValue={props.reconciledDate}/>
            <input className="bg-slate-300" type="hidden" name="transactionID"
                   defaultValue={props.transactionID}/>
            <button className="bg-blue-500 m-1 p-2 font-bold"
                    type="submit">Reconcile
            </button>
        </form>
    )
}