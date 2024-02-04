'use client'
import styles from '@/app/ui//Home.module.css';
import useSWR from 'swr'
const fetcher = (...args) => fetch(...args).then((res) => res.json())

console.log(process.env.NEXT_PUBLIC_MIMIRLEDGER_API_URL);
const myURL = new URL('/accounts', process.env.NEXT_PUBLIC_MIMIRLEDGER_API_URL);
console.log(myURL);

/*
async function getData() {
    console.log(process.env.NEXT_PUBLIC_MIMIRLEDGER_API_URL);

   const res = await fetch(myURL)
    const data = await res.json();
    console.log(data);

    // The return value is *not* serialized
    // You can return Date, Map, Set, etc.
    if (!res.ok) {
        // This will activate the closest `error.js` Error Boundary
        throw new Error('Failed to fetch data')
    }
    return res.json()
}
*/
export default function Page() {

    const { data, error, isValidating } = useSWR(myURL, fetcher)

    if (error) return <div>Failed to load</div>
    if (isValidating) return <div className="Loading">Loading...</div>;


    return (
        <div>
            <div>MY LIST OF Accounts</div>
            {data.accounts &&
                data.accounts.map((account, index) => <h1 key={index}>{account.account_name}<span>{account.account_parent}</span></h1>)}
        </div>
    );
}
