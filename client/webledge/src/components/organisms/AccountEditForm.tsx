import { useParams } from "react-router-dom";
import { TransactionAccountPostRequest } from '../../lib/definitions';
import React, {FormEvent} from "react";
import AccountSelector from "../molecules/AccountSelector";
import AccountTypeSelector from "../molecules/AccountTypeSelector";
import {useGetAccount} from "../../lib/data";
const postFormData = async (formData: FormData) => {
    try {
        // Do a bit of work to convert the entries to a plain JS object
        const formEntries = Object.fromEntries(formData);
        const accountID = Number(formEntries.accountID)
        const myURL = new URL('/accounts/'+accountID, process.env.REACT_APP_MIMIRLEDGER_API_URL);

        const newAccount : TransactionAccountPostRequest = {
            accountParent : Number(formEntries.accountParent),
            accountName : String(formEntries.accountName),
            accountType : String(formEntries.accountType),
            accountMemo : String(formEntries.accountMemo),
        };
        var json = JSON.stringify(newAccount);
        const settings :RequestInit = {
            method: 'PUT',
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

export default function AccountEditForm(){
    const { accountID } = useParams();
    const { data, isLoading, error } = useGetAccount(accountID);

    if (isLoading) return <div className="Loading">Loading...</div>
    if (error) return <div>Failed to load</div>

     return (
        <div className="flex">
         <form className="flex" onSubmit={handleSubmit}>
            <label className="my-4 text-xl font-bold mx-4 bg-slate-200">AccountName:
                <input className="bg-slate-300 font-normal" type="text" name="accountName" defaultValue={data?.accountName}/>
            </label>
            <label className="my-4 text-xl font-bold mx-4 bg-slate-200">
                AccountParent:
                <AccountSelector name={"accountParent"} id={data?.accountParent} includeTop={true} excludeID={0}
                     multiple={false}
                     multiSize={1}/>
            </label>
            <label className="my-4 text-xl font-bold">AccountType:
                <AccountTypeSelector selectedName={data?.accountType} />
            </label>
            <label className="my-4 text-xl font-bold mx-4 bg-slate-200">Memo:
                <input className=" bg-slate-300" type="text" name="accountMemo" defaultValue={data?.accountMemo}/>
            </label>
             <div className="bg-slate-300 flex">
                 <input className=" bg-slate-300" type="hidden" name="accountID" defaultValue={data?.accountID}/>
                 <button className="p-3 font-bold" type="submit">Update</button>
             </div>
         </form>
        </div>
     );
}
