
import React, {FormEvent} from "react";
import Modal from "../molecules/Modal.tsx";

const postFormData = async () => {
    try {
        const myURL = new URL('/reports/restore', import.meta.env.VITE_APP_MIMIRLEDGER_API_URL);
        const settings :RequestInit = {
            method: 'POST',
        };
        return await fetch(myURL, settings);
    } catch (error) {
        console.error('Error making POST request:', error);
    }
}

export default function ReportRestoreForm() {
    const [showModal, setShowModal] = React.useState(false);
    const [modalBody, setModalBody] = React.useState("");
    const [modalTitle, setModalTitle] = React.useState("");

    async function handleSubmit(event: FormEvent<HTMLFormElement>) {
        event.preventDefault()
        const myResponse = await postFormData();
        if (myResponse?.status == 200) {
            window.location.reload();
        } else {
            setModalTitle("ERROR RestoringReport")
            const errData = await myResponse?.json()

            setModalBody("<h1>"+errData.err+"</h1>")
            setShowModal(true)
        }
    }
    return (
        <div className="flex w-full flex-col p-4 bg-slate-200">
            <form className="mb-2" onSubmit={handleSubmit}>
                <div className="my-4 text-xl font-bold bg-slate-200 flex flex-row">
                    <div className="w-64 mr-2 text-right"/>
                    <button className="p-3 font-bold bg-red-600 text-white" type="submit">Restore Reports</button>
                </div>
            </form>
            {modalBody}
            <Modal showModal={showModal} setShowModal={setShowModal} title={modalTitle} body={modalBody} onClose={()=>{}}/>
        </div>
    )
}
