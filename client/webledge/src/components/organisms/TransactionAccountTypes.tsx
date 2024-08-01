import type { TransactionAccountType} from  "../../lib/definitions"
import {useGetTransactionAccountTypes} from "../../lib/data";

export default function TransactionAccountTypes() {
    const { data, error, isLoading } = useGetTransactionAccountTypes()

    if (error) return <div>Failed to load</div>
    if (isLoading) return <div className="Loading">Loading...</div>;

    return (
        <div className="flex w-full flex-col md:col-span-4">
            <div className="flex grow flex-col justify-between rounded-xl bg-slate-100 p-4">
                <div className="flex">
                    <div className="w-48 font-bold">
                        Account Type
                    </div>
                    <div className="w-48 font-bold">
                        Account Sign
                    </div>
                </div>
            {data?.accountTypes && data.accountTypes.map((accountType: TransactionAccountType, index: number) => {
                return (
                    <div key={index} className="flex">
                        <div className="w-48">
                            <h1 key={index}>{accountType.name}</h1>
                        </div>
                        <div className="w-48">
                            <h1 key={index}>{accountType.sign}</h1>
                        </div>
                    </div>
                );
            })}
            </div>
        </div>
    );
}
