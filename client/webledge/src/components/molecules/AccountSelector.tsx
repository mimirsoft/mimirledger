import {Account} from "../../lib/definitions";
import React from "react";
import {useGetAccounts} from "../../lib/data";

const AccountSelector = ( props:{name:string; id:number|undefined; excludeID:number|undefined; includeTop:boolean} ) => {
    const { data, error, isLoading } = useGetAccounts()
    if (isLoading) return <div className="Loading">Loading...</div>
    if (error) return <div>Failed to load</div>

   return (
        <select name={props.name} defaultValue={props.id} className="font-normal">
            { props.includeTop && <option value="0">Top Level</option>}
            {data?.accounts && data.accounts.map((account: Account, index: number) => {
                if (account.accountID == props.excludeID) {
                    return
                }
                return (
                <option key={index} value={account.accountID}> {account.accountFullName}</option>
                );
            })
        }
        </select>
   );
};

export default AccountSelector