import { SWRConfig } from 'swr'
import { BrowserRouter, Routes, Route } from "react-router-dom";
import './App.css';
import OuterContainer from './components/templates/OuterContainer';

import HomePage from "./pages/HomePage";
import NoPage from "./pages/NoPage";

import AccountsPage from "./pages/AccountsPage";
import AccountTypesPage from "./pages/AccountTypesPage";
import TransactionsAccountPage from "./pages/TransactionsAccountPage.tsx";
import AccountEditPage from "./pages/AccountEditPage";
import AccountReconcilePage from "./pages/AccountReconcilePage";
import TransactionEditPage from "./pages/TransactionEditPage";
import ReportsPage from "./pages/ReportsPage";
import ReportsListPage from "./pages/ReportsListPage";
import ReportEditPage from "./pages/ReportEditPage";
import ReportRunPage from "./pages/ReportRunPage";

const fetcher = (...args: Parameters<typeof fetch>) => fetch(...args).then(res => res.json())

function App() {
  return (
    <SWRConfig value={{ fetcher }}>
      <BrowserRouter>
        <Routes>
          <Route path="/reports/:reportID" element={<ReportRunPage/>} />
          <Route path="/" element={<OuterContainer />}>
            <Route index element={<HomePage />} />
            <Route path="accounts" element={<AccountsPage />} />
            <Route path="accounts/:accountID" element={<AccountEditPage/>} />
            <Route path="accounttypes" element={<AccountTypesPage />} />
            <Route path="reconcile/:accountID" element={<AccountReconcilePage/>} />
            <Route path="reports" element={<ReportsPage/>} />
            <Route path="reports/list" element={<ReportsListPage/>} />
            <Route path="reports/edit/:reportID" element={<ReportEditPage/>} />
            <Route path="transactions/:transactionID" element={<TransactionEditPage/>} />
            <Route path="transactions/account/:accountID" element={<TransactionsAccountPage />} />
            <Route path="*" element={<NoPage />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </SWRConfig>
  );
}

export default App;
