import {useParams, useSearchParams, useNavigate} from "react-router-dom";
import {
    TransactionDebitCreditRequest,
    TransactionDebitCreditResponse,
    TransactionEditPostRequest
} from '../../lib/definitions';
import {FormEvent, MouseEvent} from "react";
import DebitsCreditsColumn from "../molecules/DebitsCreditsColumn";
import {useGetTransaction} from "../../lib/data";
import {parseCurrency} from "../../lib/utils";
const postFormData = async (formData: FormData, debitsCount: number, creditsCount: number) => {
    try {
        // Do a bit of work to convert the entries to a plain JS object
        const formEntries = Object.fromEntries(formData);
        const transactionID = Number(formEntries.transactionID)
        const dStr = String(formEntries.transactionDate)
        const txnDate: Date = new Date(dStr);

        console.log("debitscount"+debitsCount);
        console.log("creditsCount"+creditsCount);

        const dcSet: Array<TransactionDebitCreditRequest> = []
        //iterate over all the debit/credit inputs
        // get debits
        for (let step = 0; step < debitsCount; step++) {
            let tdcAmount = parseCurrency(formEntries['debitAmount' +step])
            const tdcAccount = Number(formEntries['debitAccount' +step])
            let tdcSign = "DEBIT"
            if (tdcAmount < 0){
                tdcAmount = -tdcAmount
                tdcSign = "CREDIT"
            }
            const tdc = {accountID: tdcAccount, transactionDCAmount: tdcAmount, debitOrCredit:  tdcSign}
            console.log(tdc)
            dcSet.push(tdc)
        }
        // get credits
        for (let step = 0; step < creditsCount; step++) {
            let tdcAmount = parseCurrency(formEntries['creditAmount' +step])
            const tdcAccount = Number(formEntries['creditAccount' +step])
            let tdcSign = "CREDIT"
            if (tdcAmount < 0){
                tdcAmount = -tdcAmount
                tdcSign = "DEBIT"
            }
            const tdc = {accountID: tdcAccount, transactionDCAmount: tdcAmount, debitOrCredit: tdcSign}
            console.log(tdc)
            dcSet.push(tdc)
        }

        const myURL = new URL('/transactions/'+transactionID, import.meta.env.VITE_APP_MIMIRLEDGER_API_URL);

        const editTransaction : TransactionEditPostRequest = {
            transactionID : transactionID,
            transactionDate: txnDate.toISOString(),
            transactionComment : String(formEntries.transactionComment),
            debitCreditSet : dcSet,
        };
        const json = JSON.stringify(editTransaction);

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
    const navigate = useNavigate();
    const { transactionID } = useParams();
    const { data, isLoading, error ,
    mutate} = useGetTransaction(transactionID);
    const [searchParams] = useSearchParams();
    const returnAccountID = searchParams.get("returnAccount");

    let debitsCount = 0
    let creditsCount = 0

    async function handleSubmit(event: FormEvent<HTMLFormElement>) {
        event.preventDefault()
        const formData = new FormData(event.currentTarget)
        const response = await postFormData(formData, debitsCount, creditsCount);
        if (response?.status == 200){
            await mutate()
            navigate("/transactions/account/"+returnAccountID);
        }
        else {
            console.log("ERROR"+response)
        }
    }
    async function deleteTransaction(event:  MouseEvent<HTMLButtonElement>) {
        event.preventDefault()
        const myURL = new URL('/transactions/' + transactionID, import.meta.env.VITE_APP_MIMIRLEDGER_API_URL);
        const settings: RequestInit = {
            method: 'DELETE',
        };
        const response = await fetch(myURL, settings);
        if (response?.status == 200) {
            navigate("/transactions/account/" + returnAccountID);
        } else {
            console.log("ERROR" + response)
        }
    }

    const initialCredits: Array<TransactionDebitCreditResponse> = []
    const initialDebits: Array<TransactionDebitCreditResponse> = []

    if (isLoading) return <div className="Loading">Loading...</div>
    if (error) return <div>Failed to load</div>

    // sort debitCreditSet into debits and credits
    {data?.debitCreditSet && data.debitCreditSet.map((transaction: TransactionDebitCreditResponse) => {
        if (transaction.debitOrCredit == "CREDIT"){
            initialCredits.push(transaction)
            creditsCount++
        }
        if (transaction.debitOrCredit == "DEBIT"){
            initialDebits.push(transaction)
            debitsCount++
        }
    })}
    const addDebitsCount = (change:number) =>{
        debitsCount=debitsCount+change
    }
    const addCreditsCount = (change:number) =>{
        creditsCount=creditsCount+change
    }
    console.log(data)
    // render debits and credits
    const txnDate: Date = new Date(String(data?.transactionDate))
     return (
         <div className="flex w-full flex-col md:col-span-4 grow justify-between rounded-xl bg-slate-100 p-4">
             <div className="text-xl font-bold">
                 Edit Transaction
             </div>
             <form onSubmit={handleSubmit}>
                 <div className="flex flex-row flex-wrap">
                     <div className="my-2 mx-2 flex flex-col flex-wrap">
                         <label className="my-2 text-xl font-bold bg-slate-200">Date:
                         </label>
                         <input className="bg-slate-300 text-xl font-normal" type="date" name="transactionDate"
                                defaultValue={txnDate.toISOString().split('T')[0]}/>

                     </div>
                     <div className="my-2 mx-2 flex flex-col flex-wrap">
                         <label className="my-2 w-80 text-xl font-bold bg-slate-200">Comment
                         </label>
                         <input className="w-80 text-xl bg-slate-300 font-normal" type="text"
                                name="transactionComment"
                                defaultValue={data?.transactionComment}/>
                     </div>
                     <DebitsCreditsColumn name="debit"
                                          transactionID={data?.transactionID}
                                          dcSet={initialDebits}
                                          addCount={addDebitsCount}/>
                     <DebitsCreditsColumn name="credit"
                                          transactionID={data?.transactionID}
                                          dcSet={initialCredits}
                                          addCount={addCreditsCount}/>
                     <div className="flex my-2 mr-2">
                         <input className=" bg-slate-300" type="hidden" name="transactionID"
                                defaultValue={data?.transactionID}/>
                         <button className="bg-slate-300 my-2 p-3 font-bold" type="submit">Update</button>
                     </div>
                    <div className="flex my-2">
                    <button onClick={deleteTransaction}  className="bg-slate-300 my-2 p-3 font-bold text-red-500" type="submit">Delete</button>
                     </div>
                 </div>
             </form>
         </div>
     );
}
