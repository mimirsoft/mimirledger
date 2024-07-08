import React, {FormEvent, ReactElement} from "react";
import {ReportPostRequest, ReportBody, Report} from "../../lib/definitions";
import { useGetReports} from "../../lib/data";
import Modal from "../molecules/Modal";
import ReportSelector from "../molecules/ReportSelector";
import AccountSelector from "../molecules/AccountSelector";

const postFormData = async (formData: FormData) => {
    try {
        const myURL = new URL('/reports', process.env.REACT_APP_MIMIRLEDGER_API_URL);
        // Do a bit of work to convert the entries to a plain JS object
        const formEntries = Object.fromEntries(formData);
        console.log(formEntries.accountParent)

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
        var json = JSON.stringify(newReport);
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
type ErrResponse = {
    statusCode: number
    err: string
}
export default function ReportsRunForm(){
    const [showModal, setShowModal] = React.useState(false);
    const [modalBody, setModalBody] = React.useState("");
    const [modalTitle, setModalTitle] = React.useState("");

    async function handleSubmit(event: FormEvent<HTMLFormElement>) {
        // open report render in new window
        event.preventDefault()
        const formData = new FormData(event.currentTarget)
        const myResponse = await postFormData(formData);
        if (myResponse?.status == 200) {
            window.location.reload();
        } else {
            setModalTitle("ERROR SAVING REPORT")
            let errMsg= ""
            let errData = await myResponse?.json()

            setModalBody("<h1>"+errData.err+"</h1>")
            setShowModal(true)
        }
    };

    return (
    <div>
        <div className="flex w-full flex-col md:col-span-4 bg-slate-100 p-4">
            <div className="text-xl font-bold">
                Run Report
            </div>
            <form className="mb-2" onSubmit={handleSubmit}>
                <div className="flex-col w-fit">
                    <div className="my-4 mr-4 text-xl font-bold bg-slate-200">Report Name:
                        <input className="bg-slate-300 font-normal" type="text" name="reportName"/>
                    </div>
                    <div className="my-4 mr-4 text-xl font-bold bg-slate-200">On Accounts:
                        <AccountSelector name={"userSuppliedAccounts"} id={0}
                                         includeTop={false}
                                         excludeID={0}
                                         multiple={true}
                                         multiSize={10}/>
                    </div>
                    <div className=" flex">
                        <button className="p-3 font-bold bg-slate-300" type="submit">Run Report</button>
                    </div>
                </div>
            </form>
        </div>
        {modalBody}
        <Modal showModal={showModal} setShowModal={setShowModal} title={modalTitle} body={modalBody}/>
    </div>
    );
}
