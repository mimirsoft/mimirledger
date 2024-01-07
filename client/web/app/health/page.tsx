"use client";
import styles from '@/app/ui//Home.module.css';

export default function Page() {

    const callAPI = async () => {
        console.log(process.env.NEXT_PUBLIC_MIMIRLEDGER_API_URL);

        const myURL = new URL('/health', process.env.NEXT_PUBLIC_MIMIRLEDGER_API_URL);
        try {
            const res = await fetch(myURL);
            const data = await res.text();
            console.log(data);
        } catch (err) {
            console.log(err);
        }
    };
    return (
        <div className={styles.container}>
            <main className={styles.main}>
                <div className="flex h-20 shrink-0 items-end rounded-lg bg-blue-500 p-4 md:h-52">
                </div>
                <button onClick={callAPI}>Make API HEALTH call</button>
            </main>
            <div className="flex flex-col justify-center gap-6 rounded-lg bg-gray-50 px-6 py-10 md:w-2/5 md:px-20">
                <div className={styles.shape}></div>
            </div>
        </div>

    );
}
