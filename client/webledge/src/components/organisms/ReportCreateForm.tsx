import React, {FormEvent} from "react";
import {ReportPostRequest, ReportBody} from "../../lib/definitions";
import ReportsList from "../molecules/ReportsList"
import Modal from "../molecules/Modal";

const postFormData = async (formData: FormData) => {
    try {
        const myURL = new URL('/reports', import.meta.env.VITE_APP_MIMIRLEDGER_API_URL);
        // Do a bit of work to convert the entries to a plain JS object
        const formEntries = Object.fromEntries(formData);
        console.log(formEntries.accountParent)

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

export default function ReportCreateForm(){
    const [showModal, setShowModal] = React.useState(false);
    const [modalBody, setModalBody] = React.useState("");
    const [modalTitle, setModalTitle] = React.useState("");


    async function handleSubmit(event: FormEvent<HTMLFormElement>) {
        event.preventDefault()
        const formData = new FormData(event.currentTarget)
        const myResponse = await postFormData(formData);
        if (myResponse?.status == 200) {
            window.location.reload();
        } else {
            setModalTitle("ERROR SAVING REPORT")
            const errData = await myResponse?.json()

            setModalBody("<h1>"+errData.err+"</h1>")
            setShowModal(true)
        }
    }

    return (
    <div>
        <div className="flex w-full flex-col md:col-span-4 bg-slate-100 p-4">
            <div className="text-xl font-bold">
                Create New Report
            </div>
            <form className="mb-2" onSubmit={handleSubmit}>
                <div className="flex-col w-fit">
                    <div className="my-4 text-xl font-bold bg-slate-200 flex flex-row">
                        <div className="w-64 mr-2 text-right">ReportName:</div>
                        <input className="bg-slate-300 font-normal" type="text" name="reportName"/>
                    </div>
                    <div className="my-4 text-xl font-bold bg-slate-200 flex flex-row">
                        <div className="w-64 mr-2 text-right">Source Account Set Type:</div>
                        <select name="sourceAccountSetType" className="font-normal">
                            <option value="NONE">NONE</option>
                            <option value="GROUP">GROUP</option>
                            <option value="PREDEFINED">PREDEFINED</option>
                            <option value="USER_SUPPLIED">USER_SUPPLIED</option>
                        </select>
                    </div>
                    <div className="my-4 text-xl font-bold bg-slate-200 flex flex-row">
                        <div className="w-64 mr-2 text-right">Source Account Group:</div>
                        <select name="sourceAccountGroup" className="font-normal">
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
                        <select name="filterAccountSetType" className="font-normal">
                            <option value="NONE">NONE</option>
                            <option value="GROUP">GROUP</option>
                            <option value="PREDEFINED">PREDEFINED</option>
                            <option value="USER_SUPPLIED">USER_SUPPLIED</option>
                        </select>
                    </div>
                    <div className="my-4 text-xl font-bold bg-slate-200 flex flex-row">
                        <div className="w-64 mr-2 text-right">Filter Account Group:</div>
                        <select name="filterAccountGroup" className="font-normal">
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
                        <select name="dataSetType">
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
                    <div className="my-4 text-xl font-bold bg-slate-200 flex flex-row">
                        <div className="w-64 mr-2 text-right"/>
                        <button className="p-3 font-bold bg-blue-500 text-white" type="submit">Create Report</button>
                    </div>
                </div>
            </form>

        </div>
        <ReportsList/>
        {modalBody}
        <Modal showModal={showModal} setShowModal={setShowModal} title={modalTitle} body={modalBody} onClose={()=>{}}/>
    </div>
    );
}
