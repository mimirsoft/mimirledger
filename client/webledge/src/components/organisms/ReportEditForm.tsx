import { useParams } from "react-router-dom";
import {ReportBody, ReportPostRequest} from '../../lib/definitions';
import React, {FormEvent} from "react";
import {useGetReport} from "../../lib/data";
const postFormData = async (formData: FormData) => {
    try {
        // Do a bit of work to convert the entries to a plain JS object
        const formEntries = Object.fromEntries(formData);
        const reportID = Number(formEntries.reportID)
        const myURL = new URL('/reports/'+reportID, import.meta.env.VITE_APP_MIMIRLEDGER_API_URL);
        const newReportBody : ReportBody = {
            accountSetType:  String(formEntries.accountSetType),
            accountGroup:  String(formEntries.accountGroup),
            predefinedAccounts: [],
            recurseSubAccounts: 0,
            dataSetType:  String(formEntries.dataSetType),
        }
        const newReport : ReportPostRequest = {
            reportName : String(formEntries.reportName),
            reportBody:newReportBody,
        };
        const json = JSON.stringify(newReport);
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
    if (result?.status == 200) {
        window.location.reload();
    }
    else {
        console.log("ERROR"+result)
    }

}

export default function ReportEditForm(){
    const { reportID } = useParams();
    const { data, isLoading, error } = useGetReport(reportID);

    if (isLoading) return <div className="Loading">Loading...</div>
    if (error) return <div>Failed to load</div>

     return (
         <div className="flex w-full flex-col p-4 ">
             <div className="text-xl font-bold mb-2">
                 Edit Report
             </div>
             <form className="flex-col w-fit" onSubmit={handleSubmit}>
                 <label className="my-4 text-xl font-bold mx-4 bg-slate-200">ReportName:
                     <input className="bg-slate-300 font-normal" type="text" name="reportName"
                            defaultValue={data?.reportName}/>
                 </label>
                 <div className="my-4 mr-4 text-xl font-bold bg-slate-200">Account Set Type:
                     <select name="accountSetType" defaultValue={data?.reportBody.accountSetType}
                             className="font-normal">
                         <option value="GROUP">GROUP</option>
                         <option value="PREDEFINED">PREDEFINED</option>
                         <option value="USER_SUPPLIED">USER_SUPPLIED</option>
                     </select>
                 </div>
                 <div className="my-4 mr-4 text-xl font-bold bg-slate-200">Account Group:
                     <select name="accountGroup" defaultValue={data?.reportBody.accountGroup} className="font-normal">
                         <option value="null">NULL</option>
                         <option value="ASSET">ASSET</option>
                         <option value="LIABILITY">LIABILITY</option>
                         <option value="EQUITY">EQUITY</option>
                         <option value="INCOME">INCOME</option>
                         <option value="EXPENSE">EXPENSE</option>
                     </select>
                 </div>
                 <div className="my-4 mr-4 text-xl font-bold bg-slate-200">Data Set Type:
                     <select name="dataSetType" defaultValue={data?.reportBody.dataSetType}>
                         <option value="AGING">Aging Summary</option>
                         <option value="SUMGROUPLINE">Sum/Group/Line</option>
                         <option value="BALANCE">Balance</option>
                         <option value="INCOME">Income</option>
                         <option value="EXPENSE">Expense</option>
                         <option value="INCOMETRANS">Income Transactions</option>
                         <option value="EXPENSETRANS">Expense Transactions</option>
                         <option value="NETTRANS">Net Transactions</option>
                         <option value="LEDGER">Ledger</option>
                         <option value="RECONCILIATION">Reconciliation</option>
                     </select>
                 </div>
                 <div className="w-20">
                     {JSON.stringify(data?.reportBody)}
                 </div>
                 <div className="bg-slate-300 flex">
                     <input className=" bg-slate-300" type="hidden" name="reportID" defaultValue={data?.reportID}/>
                     <button className="p-3 font-bold" type="submit">Update</button>
                 </div>
             </form>
         </div>
     );
}
