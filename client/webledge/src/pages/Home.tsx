import React, { useEffect} from 'react';

export default function Home() {
    useEffect(() => {
        document.title = "MimirLedger";
    }, []);
    return (
        <div className="text-xl">
            <h1 className="text-2xl font-bold">MimirLedger</h1>
            Double Entry Bookkeeping Ledger
        </div>
    );
}