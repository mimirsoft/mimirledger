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
        window.open("/reports/"+reportID+"?startDate="+startDate,"_self");
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
                <div className="flex-col w-fit">
                    <div className="my-4 mr-4 text-xl font-bold bg-slate-200">Report Name:
                        <ReportSelector name={"reportID"} id={0} />
                    </div>
                    <div className="my-4 mx-2 text-xl font-bold bg-slate-200">Start Date:
                        <input className="bg-slate-300 text-xl font-normal" type="date" name="startDate"
                        defaultValue={lastMonthDate.toISOString().split('T')[0]}/>
                    </div>
                    <div className="my-4 mx-2 text-xl font-bold bg-slate-200">End Date:
                        <input className="bg-slate-300 text-xl font-normal" type="date" name="startDate"
                               defaultValue={todayDate.toISOString().split('T')[0]}/>
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
