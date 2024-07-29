import {
    Account,
} from '../../lib/definitions';
import {Link} from "react-router-dom";
import {useGetAccounts} from "../../lib/data";
import {formatCurrency} from "../../lib/utils";



export default function TransactionAccounts(){
    const { data, error, isLoading } = useGetAccounts()

    if (isLoading) return <div className="Loading">Loading...</div>
    if (error) return <div>Failed to load</div>
    let rowColor = "bg-slate-200"
    const minDate = new Date('0001-01-01T00:00:00Z');
    minDate.setDate(minDate.getDate() + 1);
    const reconcileDate: Date = new Date();
    const reconcileDateStr = reconcileDate.toISOString().split('T')[0]

    return (
        <div className="flex w-full flex-col md:col-span-4">
            <div className="flex grow flex-col justify-between rounded-xl bg-slate-100 p-4">
                <div className="text-xl font-bold">
                    My Accounts
                </div>
                <div className="flex">
                    <div className="w-8 font-bold">
                        ID
                    </div>
                    <div className="w-8 font-bold">
                        L
                    </div>
                    <div className="w-8 font-bold">
                        R
                    </div>
                    <div className="w-80 font-bold">
                        FullName
                    </div>
                    <div className="w-32 font-bold text-right mr-4">
                        Balance
                    </div>
                    <div className="w-80 font-bold">
                        Name
                    </div>
                    <div className="w-20 font-bold">
                        Type
                    </div>
                    <div className="w-20 font-bold">
                        Sign
                    </div>
                    <div className="w-24 font-bold">
                        Rec Date
                    </div>
                </div>
                {data?.accounts && data.accounts.map((account: Account, index: number) => {
                    if (rowColor == "bg-slate-200"){
                        rowColor = "bg-slate-400"
                    } else {
                        rowColor = "bg-slate-200"
                    }
                    let textColor = ""
                    if (account.accountBalance < 0) {
                        textColor = "text-red-500"
                    }
                    const acctReconciledDate: Date = new Date(account.accountReconcileDate);

                    let acctReconciledDateStr :string
                    if (acctReconciledDate < minDate) {
                        acctReconciledDateStr = "";
                    } else {
                        acctReconciledDateStr = acctReconciledDate.toISOString().split('T')[0]
                    }
                    return (
                        <div className={'flex '+rowColor} key={index}>
                            <Link to={'/transactions/account/' + account.accountID}
                                  className={`flex mr-4 }`}>
                                <div className="w-8">
                                    {account.accountID}
                                </div>
                                <div className="w-8">
                                    {account.accountLeft}
                                </div>
                                <div className="w-8">
                                    {account.accountRight}
                                </div>
                                <div className="w-80">
                                    {account.accountFullName}
                                </div>
                                <div className={"w-32 text-right mr-4 " + textColor}>
                                    {formatCurrency(account.accountBalance)}
                                </div>
                                <div className="w-80">
                                    {account.accountName}
                                </div>
                                <div className="w-20">
                                    {account.accountType}
                                </div>
                                <div className="w-20">
                                    {account.accountSign}
                                </div>
                                <div className="w-24">
                                    {acctReconciledDateStr}
                                </div>
                            </Link>
                            <Link to={'/accounts/' + account.accountID} className={`nav__item font-bold mr-2`}>
                                Edit Account
                            </Link>
                            <Link to={{
                                pathname: '/reconcile/' + account.accountID,
                                search: '?date=' + reconcileDateStr
                                }} className={`font-bold`}>
                                Reconcile
                            </Link>
                        </div>
                    );
                })}
            </div>
        </div>
    );
}
