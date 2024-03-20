
import React, { useEffect, useState, FormEvent  } from 'react';

const postFormData = async (formData: FormData) => {
    try {
        const myURL = new URL('/accounts', process.env.REACT_APP_MIMIRLEDGER_API_URL);
        console.log(myURL)
        var json =JSON.stringify(Object.fromEntries(formData));
        console.log(json)
        const settings :RequestInit = {
            method: 'POST',
            body: json,
        };
        const response = await fetch(myURL, settings);
        const result = await response.json();
        console.log('POST request result:', result);
    } catch (error) {
        console.error('Error making POST request:', error);
    }
}

export default function NewTransactionAccountsForm() {
     async function handleSubmit(event: FormEvent<HTMLFormElement>) {
         event.preventDefault()
         const formData = new FormData(event.currentTarget)
         const result = await postFormData(formData);
     };
    /*
};
    async function onSubmit(event: FormEvent<HTMLFormElement>) {
        event.preventDefault()
        const formData = new FormData(event.currentTarget)

        const response = await fetch('/accounts', {
            method: 'POST',
            body: formData,
        })
        const data = await response.json()

    }
                /*<form onSubmit={onSubmit}>
                <input type="text" name="account_name" />
                <button type="submit" disabled={isLoading}>
                    {isLoading ? 'Loading...' : 'Submit'}
                </button>
            </form>
     */
    return (
        <div>New Transaction Account
            <form onSubmit={handleSubmit}>
                <label>AccountName:
                    <input type="text" name="accountName"/>
                </label>
                <label>AccountParent:
                    <input type="text" name="accountParent"/>
                </label>
                <label>AccountType:
                    <input type="text" name="accountType"/>
                </label>
                <button type="submit">Create Account</button>
            </form>
        </div>
    );
}
