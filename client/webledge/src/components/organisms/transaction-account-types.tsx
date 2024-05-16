import type { TransactionAccountType} from  "../../lib/definitions"
import {useGetTransactionAccountTypes} from "../../lib/data";

export default function TransactionAccountTypes() {
    const { data, error, isLoading } = useGetTransactionAccountTypes()

    if (error) return <div>Failed to load</div>
    if (isLoading) return <div className="Loading">Loading...</div>;

    return (
        <div>
            {data?.accountTypes &&  data.accountTypes.map((accountType:TransactionAccountType, index:number) => {
                    return (
                        <div key={index} className="flex">
                            <div className="w-64">
                            <h1 key={index}>{accountType.name}</h1>
                            </div>
                            <div className="w-64">
                            <h1 key={index}>{accountType.sign}</h1>
                            </div>
                        </div>
                    );
            })}
        </div>
    );
}
