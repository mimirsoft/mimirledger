import React, {FormEvent} from "react";
import Modal from "../molecules/Modal";
import ReportSelector from "../molecules/ReportSelector";
import AccountSelector from "../molecules/AccountSelector";

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
        const formEntries = Object.fromEntries(formData);
        let reportID = Number(formEntries.reportID)
        let startDate =  String(formEntries.startDate)
        let endDate =  String(formEntries.endDate)
        window.open("/reports/"+reportID+"?startDate="+startDate+"&endDate="+endDate,"_blank");
    };
    let todayDate: Date = new Date()
    let lastMonthDate: Date = new Date()
    lastMonthDate.setMonth(lastMonthDate.getMonth() - 1);
    return (
    <div>
        <div className="flex w-full flex-col md:col-span-4 bg-slate-100 p-4">
            <div className="text-xl font-bold">
                Run Report
            </div>
            <form className="mb-2" onSubmit={handleSubmit}>
                <div className="flex-col w-80">
                    <div className="my-4 text-xl font-bold bg-slate-200 w-80 flex flex-row">
                        <div className="w-32 mr-2 text-right">Report:</div>
                        <ReportSelector name={"reportID"} id={0} />
                    </div>
                    <div className="my-4 text-xl font-bold bg-slate-200  w-80 flex flex-row">
                        <div className="w-32 mr-2 text-right">Start Date:</div>
                        <input className="bg-slate-300 text-xl font-normal" type="date" name="startDate"
                               defaultValue={lastMonthDate.toISOString().split('T')[0]}/>
                    </div>
                    <div className="my-4 text-xl font-bold bg-slate-200  w-80 flex flex-row">
                        <div className="w-32 mr-2 text-right">End Date:</div>
                        <input className="bg-slate-300 text-xl font-normal" type="date" name="endDate"
                               defaultValue={todayDate.toISOString().split('T')[0]}/>
                    </div>
                    <div className="my-4 text-xl font-bold bg-slate-200  w-80 flex flex-row">
                        <div className="w-32 mr-2 text-right">On Accounts:</div>
                        <AccountSelector name={"userSuppliedAccounts"} id={0}
                                         includeTop={false}
                                         excludeID={0}
                                         multiple={true}
                                         multiSize={10}/>
                    </div>
                    <div className=" flex flex-row-reverse w-80">
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
