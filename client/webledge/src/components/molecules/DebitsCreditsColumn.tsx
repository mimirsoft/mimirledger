import React from "react";
import {TransactionDebitCreditResponse} from "../../lib/definitions";
import AccountSelector from "./AccountSelector";


const DebitsCreditsColumn = ( props:{name:string;
    transactionID: number|undefined;
    dcSet:Array<TransactionDebitCreditResponse>;
	setCount: ()=>void} ) => {
    const [debitsCredits, setDebitsCredits] = React.useState(props.dcSet);

    let title = ""
    switch (props.name) {
        case "debit":
            title = "DEBIT"
            console.log("It is a DEBIT.");
            break;
        case "credit":
            title = "CREDIT"
            console.log("It is a CREDIT.");
            break;
    }

    const addDebitCredit = React.useCallback(
        () => {
                setDebitsCredits((debitsCredits) => ([
                    ...debitsCredits,
                    {transactionID: Number(props.transactionID),
                        accountID: 0,
                        debitOrCredit: props.name,
                        transactionDCAmount: 0}]));
                props.setCount();
        },
        []
    );
    return (

        <div className="my-2 mx-2 flex flex-col flex-wrap">
            <label className="my-2 w-80 text-xl font-bold bg-slate-200 w-full">{props.name}S
            </label>
            {debitsCredits.map((transaction: TransactionDebitCreditResponse,
                            index: number) => {
                return (
                    <div className="mx-0  flex flex-row flex-wrap text-right" key={index}>
                        <div className="w-16 text-right">
                            <input className="w-16 bg-slate-300 text-right" type="text"
                                   name={props.name+"Amount" + index}
                                   defaultValue={transaction.transactionDCAmount}/>
                        </div>
                        <div className="text-right">
                            <AccountSelector name={props.name+"Account" + index} id={transaction.accountID}
                                             includeTop={true}
                                             excludeID={0}/>
                        </div>
                    </div>
                );
            })}
            <div className="my-2 w-80 text-xl font-bold  w-full">
                <button onClick={addDebitCredit} className="bg-slate-300 my-2 p-3 font-bold" type="button"
                >Add {props.name}</button>
            </div>
        </div>
    );
};

export default DebitsCreditsColumn