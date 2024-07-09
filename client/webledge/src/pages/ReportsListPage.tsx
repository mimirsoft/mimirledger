import React from 'react'
import ReportsListForm from '../components/organisms/ReportsListForm';
import ReportsSubNav from "../components/organisms/ReportsSubNav";

export default function ReportsListPage() {
    return (
        <div>
            <ReportsSubNav />
            <ReportsListForm />
        </div>
    );
}
