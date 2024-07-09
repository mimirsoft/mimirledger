import React , { MouseEvent }from "react";
import {TransactionDebitCreditResponse} from "../../lib/definitions";
import AccountSelector from "./AccountSelector";
import {formatCurrency, formatCurrencyNoSign} from "../../lib/utils";


const DebitsCreditsColumn = ( props:{name:string;
    transactionID: number|undefined;
    dcSet:Array<TransactionDebitCreditResponse>;
    addCount: (change:number)=>void} ) => {
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
    const removeDebitCredit = React.useCallback((event: MouseEvent<HTMLButtonElement>) => {
        const index = parseInt(String(event?.currentTarget.dataset['index']), 10);
        setDebitsCredits((debitsCredits) => {
            const newDebitsCredits = [...debitsCredits];
            newDebitsCredits.splice(index, 1);
            return newDebitsCredits;
        });
        props.addCount(-1);
    }, []);


    const addDebitCredit = React.useCallback(() => {
        setDebitsCredits((debitsCredits) => ([
            ...debitsCredits,
            {transactionID: Number(props.transactionID),
                accountID: 0,
                debitOrCredit: props.name,
                transactionDCAmount: 0}]));
        props.addCount(1);
        },
        []
    );
    return (
        <div className="my-2 mx-2 flex flex-col flex-wrap">
            <label className="my-2 w-96 text-xl font-bold bg-slate-200 w-full">{title}S
            </label>
            {debitsCredits.map((transaction: TransactionDebitCreditResponse,
                            index: number) => {
                return (
                    <div className="mx-0 mb-2 flex flex-row flex-wrap text-right text-xl " key={index}>
                        <AccountSelector name={props.name+"Account" + index} id={transaction.accountID}
                            includeTop={true}
                            excludeID={0}
                            multiple={false}
                            multiSize={1}/>
                        <input className="w-24 text-xl bg-slate-300 text-right" type="text"
                               name={props.name+"Amount" + index}
                               defaultValue={formatCurrencyNoSign(transaction.transactionDCAmount)}/>
                        <button className="mx-1 text-xl font-bold" onClick={removeDebitCredit} data-index={index}>
                            &times;
                        </button>
                    </div>
                );
            })}
            <div className="my-0 w-96 text-xl font-bold  w-full">
                <button onClick={addDebitCredit} className="bg-slate-300 my-2 p-3 font-bold" type="button"
                >Add {title}</button>
            </div>
        </div>
    );
};

export default DebitsCreditsColumn