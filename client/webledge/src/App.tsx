import React from 'react';
import { SWRConfig } from 'swr'
import { BrowserRouter, Routes, Route } from "react-router-dom";
import './App.css';
import Home from "./pages/Home";
import AccountsPage from "./pages/AccountsPage";

import AccountTypesPage from "./pages/AccountTypesPage";
import TransactionsAccount from "./pages/TransactionsAccount";
import NoPage from "./pages/NoPage";
import OuterContainer from './components/templates/OuterContainer';
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
            <Route index element={<Home />} />
            <Route path="accounts" element={<AccountsPage />} />
            <Route path="accounts/:accountID" element={<AccountEditPage/>} />
            <Route path="accounttypes" element={<AccountTypesPage />} />
            <Route path="reconcile/:accountID" element={<AccountReconcilePage/>} />
            <Route path="reports" element={<ReportsPage/>} />
            <Route path="reports/list" element={<ReportsListPage/>} />
            <Route path="reports/edit/:reportID" element={<ReportEditPage/>} />
             <Route path="transactions/:transactionID" element={<TransactionEditPage/>} />
            <Route path="transactions/account/:accountID" element={<TransactionsAccount />} />
            <Route path="*" element={<NoPage />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </SWRConfig>
  );
}

export default App;
