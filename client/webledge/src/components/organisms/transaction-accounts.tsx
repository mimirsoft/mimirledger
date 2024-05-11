import { TransactionAccount, TransactionAccountPostRequest } from '../../lib/definitions';
import {useState} from "react";
import Modal from '../molecules/Modal'
import useSWR from "swr";
import React, {FormEvent} from "react";
import {Link} from "react-router-dom";
import AccountSelector from "../molecules/account-selector";
import AccountTypeSelector from "../molecules/account-type-selector";
const fetcher = (...args: Parameters<typeof fetch>) => fetch(...args).then(res => res.json())
const myURL = new URL('/accounts', process.env.REACT_APP_MIMIRLEDGER_API_URL);

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
    const [openModal, setOpenModal]  = useState(false)
    const { data, error, isLoading } = useSWR(myURL, fetcher)

    if (isLoading) return <div className="Loading">Loading...</div>
    if (error) return <div>Failed to load</div>

    return (
        <div className="flex w-full flex-col md:col-span-4">
            <button className="w-80 bg-blue-500 openModalbtn">Open the modal</button>
            <Modal open={openModal}></Modal>
            <h2 className={` mb-4 text-xl md:text-2xl`}>
            My Accounts
            </h2>
        <div className="flex grow flex-col justify-between rounded-xl bg-slate-100 p-4">
            <div className="text-xl font-bold">
                Create New Account
            </div>
            <div className="flex">
            <form className="flex" onSubmit={handleSubmit}>
            <label className="my-4 text-xl font-bold mx-4 bg-slate-200">AccountName:
                <input className="bg-slate-300 font-normal" type="text" name="accountName"/>
            </label>
            <label className="my-4 text-xl font-bold mx-4 bg-slate-200">
                AccountParent:
                <AccountSelector name={"accountParent"} id={0} includeTop={true} excludeID={0}/>
            </label>
            <label className="my-4 text-xl font-bold">AccountType:
                <AccountTypeSelector selectedName="" />
            </label>
            <label className="my-4 text-xl font-bold mx-4 bg-slate-200">Memo:
                <input className=" bg-slate-300" type="text" name="accountMemo"/>
            </label>
            <div className="bg-slate-300 flex">
            <button className="p-3 font-bold" type="submit">Create Account</button>
            </div>
            </form>
            </div>
        <div className="flex">
            <div className="w-8">
                ID
            </div>
            <div className="w-8">
                L
            </div>
            <div className="w-8">
                R
            </div>
            <div className="w-80">
                FullName
            </div>
            <div className="w-80">
                Name
            </div>
            <div className="w-32">
                Type
            </div>
            <div className="w-32">
                Balance
            </div>
        </div>
        {data.accounts && data.accounts.map((account: TransactionAccount, index: number) => {
            return (
                <div className="flex" key={index} >
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
                    <div className="w-80">
                        {account.accountName}
                    </div>
                    <div className="w-32">
                        {account.accountType}
                    </div>
                    <div className="w-32">
                        {account.accountBalance}
                    </div>
                    <Link to={'/transactions/account/'+account.accountID} className={`nav__item p-4 }`}>
                        LEDGER
                    </Link>
                    <Link to={'/accounts/'+account.accountID} className={`nav__item p-4 }`}>
                        EDIT ACCOUNT
                    </Link>
                </div>
            );
        })}
        </div>
    </div>
    );
}
