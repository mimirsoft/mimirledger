import React, { useEffect} from 'react';

export default function Home() {
    useEffect(() => {
        document.title = "MimirLedger";
    }, []);

    return (
        < >
            <h1>MimirLedger</h1>
        </>
    );
}