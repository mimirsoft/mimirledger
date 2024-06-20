import React from 'react'
import {useParams} from "react-router-dom";
import { useGetUnreconciledTransactionOnAccount} from "../../lib/data";
export default function AccountReconcileForm() {
    const { accountID } = useParams();
    let reconcileDate: Date = new Date();
    let reconcileDateStr = reconcileDate.toISOString().split('T')[0]
    const { data,
        isLoading,
        error } = useGetUnreconciledTransactionOnAccount(accountID, reconcileDateStr);

    if (isLoading) return <div className="Loading">Loading...</div>
    if (error) return <div>Failed to load</div>

    return <div className="flex">AccountReconcileForm</div>
}