import { useParams } from "react-router-dom";
import {ReportBody, ReportPostRequest} from '../../lib/definitions';
import React, {FormEvent} from "react";
import {useGetReport} from "../../lib/data";
const postFormData = async (formData: FormData) => {
    try {
        // Do a bit of work to convert the entries to a plain JS object
        const formEntries = Object.fromEntries(formData);
        const reportID = Number(formEntries.reportID)
        const myURL = new URL('/reports/'+reportID, process.env.REACT_APP_MIMIRLEDGER_API_URL);
        const newReportBody : ReportBody = {
            accountSetType:  String(formEntries.accountSetType),
            accountGroup:  String(formEntries.accountGroup),
            predefinedAccounts: [],
            recurseSubAccounts: 0,
        }
        const newReport : ReportPostRequest = {
            reportName : String(formEntries.reportName),
            reportBody:newReportBody,
        };
        var json = JSON.stringify(newReport);
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

export default function ReportEditForm(){
    const { reportID } = useParams();
    const { data, isLoading, error } = useGetReport(reportID);

    if (isLoading) return <div className="Loading">Loading...</div>
    if (error) return <div>Failed to load</div>

     return (
        <div className="flex">
            <form className="flex" onSubmit={handleSubmit}>
                <label className="my-4 text-xl font-bold mx-4 bg-slate-200">ReportName:
                    <input className="bg-slate-300 font-normal" type="text" name="reportName"
                           defaultValue={data?.reportName}/>
                </label>
                <div className="w-20">
                    {JSON.stringify(data?.reportBody)}
                </div>
                <div className="bg-slate-300 flex">
                    <input className=" bg-slate-300" type="hidden" name="accountID" defaultValue={data?.reportID}/>
                    <button className="p-3 font-bold" type="submit">Update</button>
                </div>
            </form>
        </div>
     );
}
