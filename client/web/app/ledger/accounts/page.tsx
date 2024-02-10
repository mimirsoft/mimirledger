'use client'
import styles from '@/app/ui//Home.module.css';

import TransactionAccounts from '@/app/ui/dashboard/transaction-account';

export default function Page() {

    async function create(formData: FormData) {
        // Logic to mutate data...
    }
    return (
        <div>
            <div>MY LIST OF Accounts</div>
            <TransactionAccounts />
            <form action={create}>...</form>
        </div>
    );
}
