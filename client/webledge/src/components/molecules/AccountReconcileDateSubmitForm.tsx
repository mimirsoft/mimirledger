import  {FormEvent} from "react";
import {AccountReconcileDatePostRequest} from "../../lib/definitions";
const updateReconcileDate= async (formData: FormData) => {
    try {
        // Do a bit of work to convert the entries to a plain JS object
        const formEntries = Object.fromEntries(formData);
        const accountID = Number(formEntries.accountID)
        const dStr = String(formEntries.reconcileDate)
        const txnDate: Date = new Date(dStr);

        const myURL = new URL('/accounts/'+accountID+"/reconciled", import.meta.env.VITE_APP_SERVER_API_URL);

        const reconciledPostRequest : AccountReconcileDatePostRequest = {
            accountID : accountID,
            accountReconcileDate: txnDate.toISOString(),
        };
        const json = JSON.stringify(reconciledPostRequest);
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

export default function AccountReconcileDateSubmitForm ( props:{
    accountID: number;
    reconciledDifferenceRemaining: number;
    reconcileDate: string;
   } )  {
    async function updateReconciledDateOnAccount(event: FormEvent<HTMLFormElement>) {
        event.preventDefault()
        const formData = new FormData(event.currentTarget)
        const response = await updateReconcileDate(formData);
        if (response?.status == 200){
            console.log("WOOT 2000")
        }
        else {
            console.log("ERROR"+response)
        }
    }


    if (props.reconciledDifferenceRemaining != 0){
        return <></>
    }
    return (
        <div className="bg-slate-300 flex mr-4">
        <form key={'reconcileDateForm' + props.accountID} onSubmit={updateReconciledDateOnAccount}>
            <input type="hidden" name="reconcileDate"
                   defaultValue={props.reconcileDate}/>
            <input  type="hidden" name="accountID"
                   defaultValue={props.accountID}/>
            <button className="w-48 bg-blue-500 font-bold text-center text-white h-16"
                    type="submit">Record Reconcile Date for Account
            </button>
        </form>
        </div>
    )
}
