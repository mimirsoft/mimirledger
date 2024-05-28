import {useParams, useSearchParams, useNavigate, Link} from "react-router-dom";
import {
    TransactionDebitCreditRequest,
    TransactionDebitCreditResponse,
    TransactionEditPostRequest
} from '../../lib/definitions';
import React, {FormEvent} from "react";
import AccountSelector from "../molecules/AccountSelector";
import DebitsCreditsColumn from "../molecules/DebitsCreditsColumn";
import {useGetTransaction} from "../../lib/data";

const postFormData = async (formData: FormData, debitsCount: number, creditsCount: number) => {
    try {
        // Do a bit of work to convert the entries to a plain JS object
        const formEntries = Object.fromEntries(formData);
        const transactionID = Number(formEntries.transactionID)
        const accountID = Number(formEntries.accountID)
        const amount = Number(formEntries.amount)
        console.log("debitscount"+debitsCount);
        console.log("creditsCount"+creditsCount);

        let dcSet: Array<TransactionDebitCreditRequest> = []
        //iterate over all the debit/credit inputs
        // get debits
        for (let step = 0; step < debitsCount; step++) {
            const tdcAmount = Number(formEntries['debitAmount' +step])
            const tdcAccount = Number(formEntries['debitAccount' +step])
            let tdc = {accountID: tdcAccount, transactionDCAmount: tdcAmount, debitOrCredit: "DEBIT" }
            console.log(tdc)
            dcSet.push(tdc)
        }
        // get credits
        for (let step = 0; step < creditsCount; step++) {
            const tdcAmount = Number(formEntries['creditAmount' +step])
            const tdcAccount = Number(formEntries['creditAccount' +step])
            let tdc = {accountID: tdcAccount, transactionDCAmount: tdcAmount, debitOrCredit: "CREDIT" }
            console.log(tdc)
            dcSet.push(tdc)
        }

        const myURL = new URL('/transactions/'+transactionID, process.env.REACT_APP_MIMIRLEDGER_API_URL);

        const editTransaction : TransactionEditPostRequest = {
            transactionID : transactionID,
            transactionComment : String(formEntries.transactionComment),
            debitCreditSet : dcSet,
        };
        var json = JSON.stringify(editTransaction);
        console.log(json);

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
    let returnAccountID = searchParams.get("returnAccount");

    let debitsCount = 0
    let creditsCount = 0

    async function handleSubmit(event: FormEvent<HTMLFormElement>) {
        event.preventDefault()
        const formData = new FormData(event.currentTarget)
        const response = await postFormData(formData, debitsCount, creditsCount);
        if (response?.status == 200){
            navigate("/transactions/account/"+returnAccountID);
        }
        else {
            console.log(response)
        }
    };

    let initialCredits: Array<TransactionDebitCreditResponse> = []
    let initialDebits: Array<TransactionDebitCreditResponse> = []

    const navigate = useNavigate();
    const { transactionID } = useParams();
    const { data, isLoading, error } = useGetTransaction(transactionID);

    if (isLoading) return <div className="Loading">Loading...</div>
    if (error) return <div>Failed to load</div>


    // sort debitCreditSet into debits and credits
    {data?.debitCreditSet && data.debitCreditSet.map((transaction: TransactionDebitCreditResponse,
                                                      index: number) => {
        if (transaction.debitOrCredit == "CREDIT"){
            initialCredits.push(transaction)
            creditsCount++
        }
        if (transaction.debitOrCredit == "DEBIT"){
            initialDebits.push(transaction)
            debitsCount++
        }
    })}
    const setDebitsCount = () =>{
        debitsCount++
    }
    const setCreditsCount = () =>{
        creditsCount++
    }
    // render debits and credits
     return (
         <div >
             <form onSubmit={handleSubmit}>
                 <div className="flex flex-row flex-wrap">
                     <div className="my-2 mx-2 flex flex-col flex-wrap">
                         <label className="my-2 w-80 text-xl font-bold bg-slate-200">Comment
                         </label>
                         <input className="w-80 bg-slate-300 font-normal" type="text"
                                name="transactionComment"
                                defaultValue={data?.transactionComment}/>
                     </div>
                     <DebitsCreditsColumn name="debit"
                                          transactionID={data?.transactionID}
                                          dcSet={initialDebits }
                                        setCount={setDebitsCount}/>
                     <DebitsCreditsColumn name="credit"
                                          transactionID={data?.transactionID}
                                          dcSet={initialCredits }
                                          setCount={setCreditsCount}/>
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
