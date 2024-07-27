import React from 'react'
import ReportCreateForm from '../components/organisms/ReportCreateForm';
import ReportsSubNav from "../components/organisms/ReportsSubNav";

export default function ReportsListPage() {
    return (
        <div>
            <ReportsSubNav />
            <ReportCreateForm />
        </div>
    );
}
