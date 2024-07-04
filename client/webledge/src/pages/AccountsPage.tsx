import React, { useEffect, useState, FormEvent } from 'react'
import TransactionAccounts from '../components/organisms/TransactionAccounts';
import AccountsSubNav from '../components/organisms/AccountsSubNav';

export default function AccountsPage() {
    return (
        <div>
            <AccountsSubNav />
            <TransactionAccounts />
        </div>
    );
}
