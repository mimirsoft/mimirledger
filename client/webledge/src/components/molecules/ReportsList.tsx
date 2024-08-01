import {Report} from "../../lib/definitions.ts";
import {Link} from "react-router-dom";
import {useGetReports} from "../../lib/data.tsx";

export default function ReportsList() {
    const {data, error, isLoading} = useGetReports()
    if (isLoading) return <div className="Loading">Loading...</div>
    if (error) return <div>Failed to load</div>

    return (
        <div className="flex w-full flex-col p-4 bg-slate-200">
            <div className="text-xl font-bold mb-2">
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
            {
                data?.reports && data.reports.map((report: Report, index: number) => {
                    return (
                        <Link to={{
                            pathname: '/reports/edit/' + report.reportID,
                        }} className={`font-bold`}>
                            <div className={'flex '} key={index}>
                                <div className="w-32">
                                    {report.reportName}
                                </div>
                                <div className="w-20 font-normal">
                                    {JSON.stringify(report.reportBody)}
                                </div>
                            </div>
                        </Link>
                    );
                })
            }
        </div>
    )
}
