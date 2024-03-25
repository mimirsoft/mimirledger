import React, { useState } from "react";
import NavButton from "../molecules/NavButton";
import NavTitle from "../molecules/NavTitle";

const MENU_LIST = [
    { text: "Home", href: "/" },
    { text: "Accounts", href: "/accounts" },
    { text: "AccountTypes", href: "/accounttypes" },
];

const NavBar = () => {
    const [navActive, setNavActive] = useState(false);
    const [activeIdx, setActiveIdx] = useState<number>(-1);

    const siteName = "Mimir Ledger";
    return (
        <div>
            <NavTitle navLogo={siteName}/>
            <nav className={`nav`}>
                <div className={`${navActive ? "active" : ""} flex bg-slate-200`}>
                    {MENU_LIST.map((menu, idx) => (
                        <div className={`w-64 bg-blue-500 p-4 text-white`}
                             onClick={() => {
                                 setActiveIdx(idx);
                                 setNavActive(false);
                             }}
                             key={menu.text}
                        >
                            <NavButton text={menu.text} href={menu.href} active={activeIdx === idx}/>
                        </div>
                    ))}
                </div>
            </nav>
        </div>
    )
}

export default NavBar;
