import TransactionAccounts from '../components/organisms/TransactionAccounts';
import AccountsSubNav from '../components/molecules/AccountsSubNav';

export default function AccountsPage() {
    return (
        <div>
            <AccountsSubNav />
            <TransactionAccounts />
        </div>
    );
}
