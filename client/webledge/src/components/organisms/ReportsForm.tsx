import React from "react";

export default function ReportsForm(){

    return (
        <div className="flex w-full flex-col md:col-span-4">
            <div className="flex grow flex-col justify-between rounded-xl bg-slate-100 p-4">
                <div className="text-xl font-bold">
                    Create New Report
                </div>
                <div className="flex">
                    <form className="flex">
                        <label className="my-4 mr-4 text-xl font-bold bg-slate-200">Report Name:
                            <input className="bg-slate-300 font-normal" type="text" name="reportName"/>
                        </label>
                        <div className="bg-slate-300 flex">
                            <button className="p-3 font-bold" type="submit">Create Report</button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    );
}
