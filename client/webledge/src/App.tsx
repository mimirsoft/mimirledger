import React from 'react';
import ReactDOM from "react-dom/client";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import logo from './logo.svg';
import './App.css';
import Home from "./pages/Home";
import Accounts from "./pages/Accounts";

import AccountTypes from "./pages/AccountTypes";
import Transactions from "./pages/Transactions";
import TransactionsAccount from "./pages/TransactionsAccount";
import NoPage from "./pages/NoPage";
import OuterContainer from './components/templates/OuterContainer';
import AccountEditPage from "./pages/AccountEdit";

function App() {
  return (
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<OuterContainer />}>
            <Route index element={<Home />} />
              <Route path="accounts" element={<Accounts />} />
              <Route path="accounts/:accountID" element={<AccountEditPage/>} />
              <Route path="accounttypes" element={<AccountTypes />} />
              <Route path="transactions/account/:accountID" element={<TransactionsAccount />} />
              <Route path="*" element={<NoPage />} />
          </Route>
        </Routes>
      </BrowserRouter>
  );
}

export default App;
