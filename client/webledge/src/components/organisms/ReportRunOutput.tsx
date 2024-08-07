
import {useParams, useSearchParams} from "react-router-dom";
import {useGetReportOutput} from "../../lib/data.tsx";

export default function ReportRunOutput() {
    const {reportID} = useParams();
    const [searchParams] = useSearchParams();
    const startDate = searchParams.get("startDate");
    const endDate = searchParams.get("endDate");
    const startDateStr = String(startDate)
    const endDateStr = String(endDate)
    const { data, isLoading, error } = useGetReportOutput(reportID,
        startDateStr,endDateStr);
    if (isLoading) return <div className="Loading">Loading...</div>
    if (error) return <div>Failed to load</div>
    console.log(data)
    return (
        <div className="flex w-full flex-col md:col-span-4 grow justify-between rounded-xl bg-slate-100 p-4">
            <div className="flex m-2">
                <div className="w-80 my-0 font-bold mx-0 bg-slate-200">
                    {data?.reportID}
                </div>
                <div className="w-80 my-0 font-bold mx-0 bg-slate-200">
                    {data?.reportName}
                </div>
                <div className="w-80 my-0 font-bold mx-0 bg-slate-200">
                    Start:{data?.startDate}
                </div>
                <div className="w-80 my-0 font-bold mx-0 bg-slate-200">
                    End:{data?.endDate}
                </div>
            </div>
        </div>
    );
}