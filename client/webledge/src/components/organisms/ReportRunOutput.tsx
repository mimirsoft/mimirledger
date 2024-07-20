import React from 'react'
import {useParams, useSearchParams} from "react-router-dom";

export default function ReportRunOutput() {
    const {reportID} = useParams();
    let [searchParams] = useSearchParams();
    let startDate = searchParams.get("startDate");
    let endDate = searchParams.get("endDate");
    let searchDateStr = String(startDate)

    return (
        <div className="flex w-full flex-col md:col-span-4 grow justify-between rounded-xl bg-slate-100 p-4">
            <div className="flex m-2">
                <div className="w-80 my-0 font-bold mx-0 bg-slate-200">
                    {reportID}
                </div>
                <div className="w-80 my-0 font-bold mx-0 bg-slate-200">
                    {searchDateStr}
                    {endDate}
                </div>
            </div>
        </div>
    );
}