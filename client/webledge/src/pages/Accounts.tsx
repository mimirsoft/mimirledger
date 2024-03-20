import styles from '@/app/ui//Home.module.css';
import React, { useEffect, useState, FormEvent } from 'react'
import TransactionAccounts from '../components/organisms/transaction-accounts';
import NewTransactionAccountsForm from '../components/organisms/new-transaction-account-form';


export default function AccountTypes() {

    /*const [isLoading, setIsLoading] = useState<boolean>(false)

    async function onSubmit(event: FormEvent<HTMLFormElement>) {
        event.preventDefault()
        const formData = new FormData(event.currentTarget)

        const response = await fetch('/accounts', {
            method: 'POST',
            body: formData,
        })
        const data = await response.json()

    }
                /*<form onSubmit={onSubmit}>
                <input type="text" name="account_name" />
                <button type="submit" disabled={isLoading}>
                    {isLoading ? 'Loading...' : 'Submit'}
                </button>
            </form>
     */
    return (
        <div>
            <TransactionAccounts />
        </div>
    );
}
