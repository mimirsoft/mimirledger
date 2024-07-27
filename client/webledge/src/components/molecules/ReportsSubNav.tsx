import {Link, useMatch, useResolvedPath} from "react-router-dom";


export default function ReportsSubNav() {
    const resolvedPath = useResolvedPath('/reports/list')
    const isActive = useMatch({path:resolvedPath.pathname, end:true})
    return (
        <>
            <div className="flex">
                <div className={` bg-blue-600 p-2 font-bold text-white`}>
                    <Link to={'/reports/list'} className={`${isActive ? "underline" : ""} p-2 hover:underline`} >Edit Reports
                    </Link>
                </div>
            </div>
        </>
    )
}
