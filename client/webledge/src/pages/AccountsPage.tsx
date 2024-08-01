import TransactionAccounts from '../components/organisms/TransactionAccounts';
import AccountsSubNav from '../components/molecules/AccountsSubNav';
import CreateAccountForm from "../components/molecules/CreateAccountForm.tsx";

export default function AccountsPage() {
    return (
        <div>
            <AccountsSubNav />
            <CreateAccountForm />
            <TransactionAccounts />
        </div>
    );
}
