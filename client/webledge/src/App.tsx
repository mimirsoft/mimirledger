import React from 'react';
import ReactDOM from "react-dom/client";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import logo from './logo.svg';
import './App.css';
import Home from "./pages/Home";
import Accounts from "./pages/Accounts";
import AccountTypes from "./pages/AccountTypes";
import NoPage from "./pages/NoPage";
import OuterContainer from './components/templates/OuterContainer';

function App() {
  return (
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<OuterContainer />}>
            <Route index element={<Home />} />
            <Route path="accounts" element={<Accounts />} />
            <Route path="accounttypes" element={<AccountTypes />} />
            <Route path="*" element={<NoPage />} />df
          </Route>
        </Routes>
      </BrowserRouter>
  );
}

export default App;
