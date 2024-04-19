import {TransactionAccount} from "../../lib/definitions";
import React from "react";
import useSWR from "swr";

const fetcher = (...args: Parameters<typeof fetch>) => fetch(...args).then(res => res.json())
const myURL = new URL('/accounts', process.env.REACT_APP_MIMIRLEDGER_API_URL);

const AccountSelector = ( props:{id:number|undefined} ) => {
    const { data, error, isLoading } = useSWR(myURL, fetcher)
    if (isLoading) return <div className="Loading">Loading...</div>
    if (error) return <div>Failed to load</div>

    return (
        <select name="accountParent" defaultValue={props.id}>
            <option value="0">Top Level</option>
            {data.accounts && data.accounts.map((account: TransactionAccount, index: number) => {
                return (
                    <option key={index} value={account.accountID}  > {account.accountFullName}</option>
                );
            })}
        </select>

)
    ;
};

export default AccountSelector