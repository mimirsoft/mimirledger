import React from 'react';
import TransactionAccountTypes from "../components/organisms/TransactionAccountTypes";
import AccountsSubNav from "../components/organisms/AccountsSubNav";

export default function AccountTypesPage() {
    return (
        <div>
            <AccountsSubNav />
            <TransactionAccountTypes />
        </div>
    );
}
