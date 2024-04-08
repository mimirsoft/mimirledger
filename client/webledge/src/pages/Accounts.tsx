import styles from '@/app/ui//Home.module.css';
import React, { useEffect, useState, FormEvent } from 'react'
import TransactionAccounts from '../components/organisms/transaction-accounts';
import NewTransactionAccountsForm from '../components/organisms/new-transaction-account-form';


export default function Accounts() {
    return (
        <div>
            <TransactionAccounts />
        </div>
    );
}
