import {useGetAccounts} from "../../lib/data.tsx";
import AccountSelector from "./AccountSelector.tsx";
import AccountTypeSelector from "./AccountTypeSelector.tsx";
import {TransactionAccountPostRequest} from "../../lib/definitions.ts";
import {FormEvent} from "react";

const postFormData = async (formData: FormData) => {
    try {
        const myURL = new URL('/accounts', import.meta.env.VITE_APP_MIMIRLEDGER_API_URL);
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
        const json = JSON.stringify(newAccount);
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
    if (result?.status == 200) {
        window.location.reload();
    }
    else {
        console.log("ERROR"+result)
    }
}

export default function CreateAccountForm() {
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
                            <input className="font-normal bg-slate-300" type="text" name="accountMemo"/>
                        </label>
                        <div className="bg-slate-300 flex">
                            <button className="p-3 font-bold bg-blue-500 text-white" type="submit">Create Account
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    )
}