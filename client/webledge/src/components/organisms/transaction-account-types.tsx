import useSWR from 'swr'
import type { TransactionAccountType} from  "../../lib/definitions"
const fetcher = (...args: Parameters<typeof fetch>) => fetch(...args).then(res => res.json())
const myURL = new URL('/accounttypes', process.env.REACT_APP_MIMIRLEDGER_API_URL);

export default function TransactionAccountTypes() {

    const { data, error, isLoading } = useSWR(myURL, fetcher)

    if (error) return <div>Failed to load</div>
    if (isLoading) return <div className="Loading">Loading...</div>;


    return (
        <div>
            {data.accountTypes &&  data.accountTypes.map((accountType:TransactionAccountType, index:number) => {
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
