import {ArrowPathIcon} from '@heroicons/react/24/outline';
import clsx from 'clsx';
import {Transaction} from '../../lib/definitions';
import {useState} from "react";
import Modal from '../molecules/Modal'

import {useTransactionAccounts} from '../../lib/data';
import useSWR from "swr";
import React, {FormEvent} from "react";

const fetcher = (...args: Parameters<typeof fetch>) => fetch(...args).then(res => res.json())
const myURL = new URL('/transactions', process.env.REACT_APP_MIMIRLEDGER_API_URL);

// get the transaction on account

export default function TransactionLedger() {
    return (  <div className="flex w-full flex-col md:col-span-4">
        <h2 className={` mb-4 text-xl md:text-2xl`}>
            Transaction - Account
        </h2>
        </div>
        );
    }
