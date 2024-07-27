import React from 'react'
import ReportsRunForm from '../components/organisms/ReportsRunForm';
import ReportsSubNav from "../components/molecules/ReportsSubNav";

export default function ReportsPage() {
    return (
        <div>
            <ReportsSubNav />
            <ReportsRunForm />
        </div>
    );
}
