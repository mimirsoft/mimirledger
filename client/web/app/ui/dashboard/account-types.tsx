'use client'
import styles from '@/app/ui//Home.module.css';
import useSWR from 'swr'
const fetcher = (...args) => fetch(...args).then((res) => res.json())

console.log(process.env.NEXT_PUBLIC_MIMIRLEDGER_API_URL);
const myURL = new URL('/accounttypes', process.env.NEXT_PUBLIC_MIMIRLEDGER_API_URL);
console.log(myURL);

export default function AccountTypes() {

    const { data, error, isValidating } = useSWR(myURL, fetcher)

    if (error) return <div>Failed to load</div>
    if (isValidating) return <div className="Loading">Loading...</div>;


    return (
        <div>
            {data.accountTypes &&
                data.accountTypes.map((accountType, index) => <h1  key={index}>{accountType.name}</h1>)}
        </div>
    );
}
