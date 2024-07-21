import {TransactionAccountType} from "../../lib/definitions";
import React from "react";

const AccountTypeSelector = ( props:{selectedName:string|undefined} ) => {
    const accountTypes: TransactionAccountType[] = [
        { "name": "ASSET", "sign": "DEBIT" },
        { "name": "LIABILITY", "sign": "CREDIT" },
        { "name": "EQUITY", "sign": "CREDIT" },
        { "name": "INCOME", "sign": "DEBIT" },
        { "name": "EXPENSE", "sign": "CREDIT" },
        { "name": "GAIN", "sign": "CREDIT" },
        { "name": "LOSS", "sign": "DEBIT" },
    ];
    return (
        <select name="accountType" defaultValue={props.selectedName}>
            {accountTypes.map((accountType: TransactionAccountType, index: number) => {
                return (
                    <option key={index} value={accountType.name}  > {accountType.name}</option>
                );
            })}
        </select>

    )
        ;
};

export default AccountTypeSelector