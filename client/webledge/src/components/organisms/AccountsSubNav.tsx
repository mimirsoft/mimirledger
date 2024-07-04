import React from 'react'
import {Link} from "react-router-dom";


export default function AccountsSubNav() {
    return (
        <>
            <div className="flex">
                <div className={` bg-blue-600 p-2 font-bold text-white`}>
                    <Link to={'/accounttypes'} >AccountTypes
                    </Link>
                </div>
            </div>
        </>
    )
}
