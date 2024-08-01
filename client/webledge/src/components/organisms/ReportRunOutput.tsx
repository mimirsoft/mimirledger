
import {useParams, useSearchParams} from "react-router-dom";

export default function ReportRunOutput() {
    const {reportID} = useParams();
    const [searchParams] = useSearchParams();
    const startDate = searchParams.get("startDate");
    const endDate = searchParams.get("endDate");
    const searchDateStr = String(startDate)

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