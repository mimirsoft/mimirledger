import React from "react";

const DebitCreditSelector = ( props:{name:string;selectedValue:string|undefined} ) => {
    return (
        <select name={props.name} defaultValue={props.selectedValue}>
            <option value="DEBIT">DEBIT</option>
            <option value="CREDIT">CREDIT</option>
        </select>
    );
};

export default DebitCreditSelector