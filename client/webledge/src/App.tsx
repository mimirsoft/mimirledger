import React from 'react';
import ReactDOM from "react-dom/client";
import { SWRConfig } from 'swr'
import { BrowserRouter, Routes, Route } from "react-router-dom";
import logo from './logo.svg';
import './App.css';
import Home from "./pages/Home";
import Accounts from "./pages/Accounts";

import AccountTypes from "./pages/AccountTypes";
import TransactionsAccount from "./pages/TransactionsAccount";
import NoPage from "./pages/NoPage";
import OuterContainer from './components/templates/OuterContainer';
import AccountEditPage from "./pages/AccountEdit";
import TransactionEditPage from "./pages/TransactionEdit";

const fetcher = (...args: Parameters<typeof fetch>) => fetch(...args).then(res => res.json())

function App() {
  return (
      <SWRConfig value={{ fetcher }}>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<OuterContainer />}>
            <Route index element={<Home />} />
              <Route path="accounts" element={<Accounts />} />
              <Route path="accounts/:accountID" element={<AccountEditPage/>} />
              <Route path="accounttypes" element={<AccountTypes />} />
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
