import {useParams, useSearchParams, useNavigate, Link} from "react-router-dom";
import {
    TransactionDebitCreditRequest,
    TransactionDebitCreditResponse,
    TransactionEditPostRequest
} from '../../lib/definitions';
import React, {FormEvent} from "react";
import AccountSelector from "../molecules/AccountSelector";
import {useGetTransaction} from "../../lib/data";

const postFormData = async (formData: FormData) => {
    try {
        // Do a bit of work to convert the entries to a plain JS object
        const formEntries = Object.fromEntries(formData);
        const transactionID = Number(formEntries.transactionID)
        const accountID = Number(formEntries.accountID)
        const amount = Number(formEntries.amount)

        let dcSet: Array<TransactionDebitCreditRequest> = []
        //iterate over all the debit/credit inputs
        // get debits
        for (let step = 0; step < 1; step++) {
            const tdcAmount = Number(formEntries['debitAmount' +step])
            const tdcAccount = Number(formEntries['debitAccount' +step])
            let tdc = {accountID: tdcAccount, transactionDCAmount: tdcAmount, debitOrCredit: "DEBIT" }
            console.log(tdc)
            dcSet.push(tdc)
        }
        // get credits
        for (let step = 0; step < 1; step++) {
            const tdcAmount = Number(formEntries['creditAmount' +step])
            const tdcAccount = Number(formEntries['creditAccount' +step])
            let tdc = {accountID: tdcAccount, transactionDCAmount: tdcAmount, debitOrCredit: "CREDIT" }
            dcSet.push(tdc)
        }

        const myURL = new URL('/transactions/'+transactionID, process.env.REACT_APP_MIMIRLEDGER_API_URL);

        const newAccount : TransactionEditPostRequest = {
            transactionID : transactionID,
            transactionComment : String(formEntries.transactionComment),
            debitCreditSet : dcSet,
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

export default function TransactionEditForm(){
    let [searchParams] = useSearchParams();
    let returnAccountID = searchParams.get("returnAccount"); // is the string "Jonathan Smith"

    async function handleSubmit(event: FormEvent<HTMLFormElement>) {
        event.preventDefault()
        const formData = new FormData(event.currentTarget)
        const response = await postFormData(formData);
        if (response?.status == 200){
            navigate("/transactions/account/"+returnAccountID);
        }
        else {
            console.log(response)
        }

    };

    const navigate = useNavigate();
    const { transactionID } = useParams();
    const { data, isLoading, error } = useGetTransaction(transactionID);

    if (isLoading) return <div className="Loading">Loading...</div>
    if (error) return <div>Failed to load</div>

    let creditSet: Array<TransactionDebitCreditResponse> = []
    let debitSet: Array<TransactionDebitCreditResponse> = []
    {data?.debitCreditSet && data.debitCreditSet.map((transaction: TransactionDebitCreditResponse,
                                                      index: number) => {
        if (transaction.debitOrCredit == "CREDIT"){
            creditSet.push(transaction)
        }
        if (transaction.debitOrCredit == "DEBIT"){
            debitSet.push(transaction)
        }

    })}
    // sort debitCreditSet into debits and credits
    // render debits and credits
     return (
         <div >
             <form onSubmit={handleSubmit}>
                 <div className="flex flex-row flex-wrap">
                     <div className="my-2 mx-4 flex flex-col flex-wrap">
                         <label className="my-2 w-80 text-xl font-bold mx-4 bg-slate-200">Comment
                         </label>
                         <input className="mx-4 w-80 bg-slate-300 font-normal" type="text"
                                name="transactionComment"
                                defaultValue={data?.transactionComment}/>
                     </div>
                     <div className="my-2 mx-4 flex flex-col flex-wrap">
                         <label className="my-2 w-80 text-xl font-bold bg-slate-200 w-full">Debits
                         </label>
                         {debitSet.map((transaction: TransactionDebitCreditResponse,
                                                                           index: number) => {
                             if (transaction.debitOrCredit == "DEBIT") {
                                 return (
                                     <div className="mx-0  flex flex-row flex-wrap text-right" key={index}>
                                         <div className="w-16 text-right">
                                             <input className="w-16 bg-slate-300 text-right" type="text"
                                                    name={"debitAmount" + index}
                                                    defaultValue={transaction.transactionDCAmount}/>
                                         </div>
                                         <div className="text-right">
                                             <AccountSelector name={"debitAccount" + index} id={transaction.accountID}
                                                              includeTop={true}
                                                              excludeID={0}/>
                                         </div>
                                     </div>
                                 );
                             }
                             return;
                         })}
                         <div className="my-2 w-80 text-xl font-bold  w-full">
                             <button className="bg-slate-300 my-2 p-3 font-bold" type="button">Add Debit</button>
                         </div>

                     </div>
                     <div className="my-2 mx-4 flex flex-col flex-wrap">
                         <label className="my-2 w-80 text-xl font-bold bg-slate-200 w-full">Credits
                         </label>
                         {creditSet.map((transaction: TransactionDebitCreditResponse,
                                                                           index: number) => {
                             if (transaction.debitOrCredit == "CREDIT") {
                                 return (
                                     <div className="mx-0  flex flex-row flex-wrap text-right" key={index}>
                                         <div className="w-16 text-right">
                                             <input className="w-16 bg-slate-300 text-right" type="text"
                                                    name={"creditAmount" + index}
                                                    defaultValue={transaction.transactionDCAmount}/>
                                         </div>
                                         <div className="text-right">
                                             <AccountSelector name={"creditAccount" + index} id={transaction.accountID}
                                                              includeTop={true}
                                                              excludeID={0}/>
                                         </div>
                                     </div>
                                 );
                             }
                             return;
                         })}
                         <div className="my-2 w-80 text-xl font-bold  w-full">
                             <button className="bg-slate-300 my-2 p-3 font-bold" type="button">Add Credit</button>
                         </div>
                     </div>
                     <div className="flex my-2">
                         <input className=" bg-slate-300" type="hidden" name="transactionID"
                                defaultValue={data?.transactionID}/>
                         <button className="bg-slate-300 my-2 p-3 font-bold" type="submit">Update</button>
                     </div>
                 </div>
             </form>
         </div>
     );
}
