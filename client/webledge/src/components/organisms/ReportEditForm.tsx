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
            sourceAccountSetType:  String(formEntries.sourceAccountSetType),
            sourceAccountGroup:  String(formEntries.sourceAccountGroup),
            sourcePredefinedAccounts: [],
            sourceRecurseSubAccounts: false,
            sourceRecurseSubAccountsDepth: 0,
            filterAccountSetType:  String(formEntries.filterAccountSetType),
            filterAccountGroup:  String(formEntries.filterAccountGroup),
            filterPredefinedAccounts: [],
            filterRecurseSubAccounts: false,
            filterRecurseSubAccountsDepth: 0,
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
         <div className="flex w-full flex-col p-4 bg-slate-200">
             <div className="text-xl font-bold mb-2">
                 Edit Report
             </div>
             <form className="flex-col w-fit" onSubmit={handleSubmit}>
                 <div className="flex-col w-fit">
                     <div className="my-4 text-xl font-bold bg-slate-200 flex flex-row">
                         <div className="w-64 mr-2 text-right">ReportName:</div>
                         <input className="bg-slate-300 font-normal" type="text" name="reportName"
                                defaultValue={data?.reportName}/>
                     </div>
                     <div className="my-4 text-xl font-bold bg-slate-200 flex flex-row">
                         <div className="w-64 mr-2 text-right">Source Account Set Type:</div>
                         <select name="sourceAccountSetType" defaultValue={data?.reportBody.sourceAccountSetType}
                                 className="font-normal">
                             <option value="NONE">NONE</option>
                             <option value="GROUP">GROUP</option>
                             <option value="PREDEFINED">PREDEFINED</option>
                             <option value="USER_SUPPLIED">USER_SUPPLIED</option>
                         </select>
                     </div>
                     <div className="my-4 text-xl font-bold bg-slate-200 flex flex-row">
                         <div className="w-64 mr-2 text-right">Source Account Group:</div>
                         <select name="sourceAccountGroup" defaultValue={data?.reportBody.sourceAccountGroup}
                                 className="font-normal">
                             <option value="null">NULL</option>
                             <option value="ASSET">ASSET</option>
                             <option value="LIABILITY">LIABILITY</option>
                             <option value="EQUITY">EQUITY</option>
                             <option value="INCOME">INCOME</option>
                             <option value="EXPENSE">EXPENSE</option>
                         </select>
                     </div>
                     <div className="my-4 text-xl font-bold bg-slate-200 flex flex-row">
                         <div className="w-64 mr-2 text-right">Filter Account Set Type:</div>
                         <select name="filterAccountSetType" defaultValue={data?.reportBody.filterAccountSetType}
                                 className="font-normal">
                             <option value="NONE">NONE</option>
                             <option value="GROUP">GROUP</option>
                             <option value="PREDEFINED">PREDEFINED</option>
                             <option value="USER_SUPPLIED">USER_SUPPLIED</option>
                         </select>
                     </div>
                     <div className="my-4 text-xl font-bold bg-slate-200 flex flex-row">
                         <div className="w-64 mr-2 text-right">Filter Account Group:</div>
                         <select name="filterAccountGroup" defaultValue={data?.reportBody.filterAccountGroup}
                                 className="font-normal">
                             <option value="null">NULL</option>
                             <option value="ASSET">ASSET</option>
                             <option value="LIABILITY">LIABILITY</option>
                             <option value="EQUITY">EQUITY</option>
                             <option value="INCOME">INCOME</option>
                             <option value="EXPENSE">EXPENSE</option>
                         </select>
                     </div>
                     <div className="my-4 text-xl font-bold bg-slate-200 flex flex-row">
                         <div className="w-64 mr-2 text-right">Data Set Type:</div>
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
                     <div className="my-4 text-xl font-bold bg-slate-200 flex">
                         <div className="w-64 min-w-64 mr-2 text-right flex flex-col">ReportJSON:</div>
                         <div className="flex flex-col font-normal">
                         {JSON.stringify(data?.reportBody, null, 2)}
                         </div>
                     </div>
                     <div className="my-4 text-xl font-bold bg-slate-200 flex flex-row">
                         <div className="w-64 mr-2 text-right"/>
                         <input className=" bg-slate-300" type="hidden" name="reportID" defaultValue={data?.reportID}/>
                         <button className="p-3 font-bold bg-blue-500 text-white" type="submit">Update</button>
                     </div>
                 </div>
             </form>
         </div>
     );
}
