import {Report} from "../../lib/definitions";
import React from "react";
import { useGetReports} from "../../lib/data";

const ReportSelector = ( props:{name:string; id:number|undefined; } ) => {
    const { data, error, isLoading } = useGetReports()
    if (isLoading) return <div className="Loading">Loading...</div>
    if (error) return <div>Failed to load</div>

    return (
        <select name={props.name} defaultValue={props.id} className="font-normal">
            {data?.reports && data.reports.map((report: Report, index: number) => {
               return (
                    <option key={index} value={report.reportID}> {report.reportName}</option>
                );
            })
            }
        </select>
    );
};

export default ReportSelector