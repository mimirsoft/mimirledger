import styles from '@/app/ui//Home.module.css';
import React, { useEffect, useState, FormEvent } from 'react'
import TransactionAccountLedger from '../components/organisms/transaction-ledger';


export default function TransactionsAccount() {
    return (
        <div>
            <TransactionAccountLedger />
        </div>
    );
}
