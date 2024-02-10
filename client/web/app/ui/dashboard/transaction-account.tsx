import { ArrowPathIcon } from '@heroicons/react/24/outline';
import clsx from 'clsx';
import Image from 'next/image';
import { lusitana } from '@/app/ui/fonts';
import { TransactionAccount } from '@/app/lib/definitions';
import { fetchTransactionAccounts } from '@/app/lib/data';
import useSWR from "swr";
export default function TransactionAccounts(){
    const { data, error, isValidating } = fetchTransactionAccounts();

    if (error) return <div>Failed to load</div>
    if (isValidating) return <div className="Loading">Loading...</div>;

    return (
        <div className="flex w-full flex-col md:col-span-4">
<h2 className={`${lusitana.className} mb-4 text-xl md:text-2xl`}>
Latest Invoices
</h2>
<div className="flex grow flex-col justify-between rounded-xl bg-gray-50 p-4">
   <div className="bg-white px-6">
   {data.accounts && data.accounts.map((account, i) => {
     return (
       <div
         key={account.id}
         className={clsx(
           'flex flex-row items-center justify-between py-4',
           {
             'border-t': i !== 0,
           },
         )}
       >
         <div className="flex items-center">
           <Image
             src={account.image_url}
             alt={`${account.name}'s profile picture`}
             className="mr-4 rounded-full"
             width={32}
             height={32}
           />
           <div className="min-w-0">
             <p className="truncate text-sm font-semibold md:text-base">
               {account.name}
             </p>
             <p className="hidden text-sm text-gray-500 sm:block">
               {account.email}
             </p>
           </div>
         </div>
         <p
           className={`${lusitana.className} truncate text-sm font-medium md:text-base`}
         >
           {account.amount}
         </p>
       </div>
     );
   })}
 </div>
<div className="flex items-center pb-2 pt-6">
<ArrowPathIcon className="h-5 w-5 text-gray-500" />
<h3 className="ml-2 text-sm text-gray-500 ">Updated just now</h3>
</div>
</div>
</div>
);
}
