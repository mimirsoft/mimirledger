"use client";
import styles from '../styles/Home.module.css';

export default function Home() {

    const callAPI = async () => {
        const myURL = new URL('/health', process.env.MIMIRLEDGER_API_URL);
        try {
            const res = await fetch(myURL);
            const data = await res.json();
            console.log(data);
        } catch (err) {
            console.log(err);
        }
    };
    return (
      <div className={styles.container}>
        <main className={styles.main}>
          <button onClick={callAPI}>Make API call</button>
        </main>
      </div>
    );
}
