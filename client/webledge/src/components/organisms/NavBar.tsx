import NavButton from "../molecules/NavButton";
import NavTitle from "../molecules/NavTitle";

const MENU_LIST = [
    { text: "Home", href: "/" },
    { text: "Accounts", href: "/accounts" },
    { text: "Reports", href: "/reports" },
];

const NavBar = () => {
    const siteName = "MimirLedger";
    return (
        <div>
            <NavTitle navLogo={siteName}/>
            <nav className={`nav`}>
                <div className={`flex bg-slate-200`}>
                    {MENU_LIST.map((menu, idx) => (
                        <div className={`bg-blue-500 p-4 font-bold text-white text-xl`}
                            key={menu.text+idx} >
                            <NavButton text={menu.text} href={menu.href}/>
                        </div>
                    ))}
                </div>
            </nav>
        </div>
    )
}

export default NavBar;
