import React, {FormEvent, ReactElement} from "react";
import {ReportPostRequest, ReportBody, Report} from "../../lib/definitions";
import { useGetReports} from "../../lib/data";
import {Link} from "react-router-dom";
import Modal from "../molecules/Modal";

const postFormData = async (formData: FormData) => {
    try {
        const myURL = new URL('/reports', import.meta.env.VITE_APP_MIMIRLEDGER_API_URL);
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
type ErrResponse = {
    statusCode: number
    err: string
}
export default function ReportsForm(){
    const { data, error, isLoading } = useGetReports()

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
            const errMsg= ""
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
                    <div className="my-4 mr-4 text-xl font-bold bg-slate-200">Report Name:
                        <input className="bg-slate-300 font-normal" type="text" name="reportName"/>
                    </div>
                    <div className="my-4 mr-4 text-xl font-bold bg-slate-200">Account Set Type:
                        <select name="accountSetType" className="font-normal">
                            <option value="GROUP">GROUP</option>
                            <option value="PREDEFINED">PREDEFINED</option>
                            <option value="USER_SUPPLIED">USER_SUPPLIED</option>
                        </select>
                    </div>
                    <div className="my-4 mr-4 text-xl font-bold bg-slate-200">Account Group:
                        <select name="accountGroup" className="font-normal">
                            <option value="null">NULL</option>
                            <option value="ASSET">ASSET</option>
                            <option value="LIABILITY">LIABILITY</option>
                            <option value="EQUITY">EQUITY</option>
                            <option value="INCOME">INCOME</option>
                            <option value="EXPENSE">EXPENSE</option>
                        </select>
                    </div>
                    <div className="my-4 mr-4 text-xl font-bold bg-slate-200">Data Set Type:
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
                    <div className=" flex">
                        <button className="p-3 font-bold bg-slate-300" type="submit">Create Report</button>
                    </div>
                </div>
            </form>
            <div className="text-xl font-bold">
                My Reports
            </div>
            <div className="flex">
                <div className="w-32 font-bold">
                    Name
                </div>
                <div className="w-20 font-bold">
                    Body
                </div>
            </div>
            {data?.reports && data.reports.map((report: Report, index: number) => {
                return (
                    <Link to={{
                        pathname: '/reports/edit/' + report.reportID,
                    }} className={`font-bold`}>
                        <div className={'flex '} key={index}>
                            <div className="w-32">
                                {report.reportName}
                            </div>
                            <div className="w-20">
                            {JSON.stringify(report.reportBody)}
                        </div>
                    </div>
                    </Link>
                );
            })}
        </div>
        {modalBody}
        <Modal showModal={showModal} setShowModal={setShowModal} title={modalTitle} body={modalBody}/>
    </div>
    );
}
