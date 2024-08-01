import ReportCreateForm from '../components/organisms/ReportCreateForm';
import ReportsSubNav from "../components/molecules/ReportsSubNav";

export default function ReportsListPage() {
    return (
        <div>
            <ReportsSubNav />
            <ReportCreateForm />
        </div>
    );
}
