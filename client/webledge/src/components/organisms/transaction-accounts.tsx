import { ArrowPathIcon } from '@heroicons/react/24/outline';
import clsx from 'clsx';
import { TransactionAccount } from '../../lib/definitions';
import { useTransactionAccounts } from '../../lib/data';
import useSWR from "swr";
const fetcher = (...args: Parameters<typeof fetch>) => fetch(...args).then(res => res.json())
const myURL = new URL('/accounts', process.env.REACT_APP_MIMIRLEDGER_API_URL);
console.log(myURL);

export default function TransactionAccounts(){

    const { data, error, isLoading } = useSWR(myURL, fetcher)

    if (isLoading) return <div className="Loading">Loading...</div>
    if (error) return <div>Failed to load</div>

    return (
        <div className="flex w-full flex-col md:col-span-4">
<h2 className={` mb-4 text-xl md:text-2xl`}>
My Accounts
</h2>
<div className="flex grow flex-col justify-between rounded-xl bg-gray-50 p-4">
   <div className="bg-white px-6">
   {data.accounts && data.accounts.map((account:TransactionAccount, index:number) => {
     return (
       <div>
         <div className="flex items-center">
           <div className="min-w-0">

             <p className="hidden text-sm text-gray-500 sm:block">
               {account.account_fullname}
             </p>
           </div>
         </div>
         <p
           className={` truncate text-sm font-medium md:text-base`}
         >
           {account.account_balance}
         </p>
       </div>
     );
   })}
 </div>
</div>
</div>
);
}
