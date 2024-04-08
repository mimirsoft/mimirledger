import styles from '@/app/ui//Home.module.css';
import React, { useEffect, useState, FormEvent } from 'react'
import TransactionLedger from '../components/organisms/transaction-ledger';


export default function Transactions() {
    return (
        <div>
            <TransactionLedger />
        </div>
    );
}
