export const formatCurrency = (amount: number) => {
    return (amount / 100).toLocaleString('en-US', {
        style: 'currency',
        currency: 'USD',
    });
};

export const formatCurrencyNoSign = (amount: number) => {
    return (amount / 100).toLocaleString('en-US',{ minimumFractionDigits: 2 } );
};

export const parseCurrency = (formAmt:  FormDataEntryValue ) => {
    let amount = String(formAmt)
    amount = amount.replace(/,/g,"")
    amount = amount.replace("$","")
    return Number(amount)*100
};