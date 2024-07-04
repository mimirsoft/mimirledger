import React, { useState } from "react";
import NavButton from "../molecules/NavButton";
import NavTitle from "../molecules/NavTitle";

const MENU_LIST = [
    { text: "Home", href: "/" },
    { text: "Accounts", href: "/accounts" },
    { text: "Reports", href: "/reports" },
];

const NavBar = () => {
    const [navActive, setNavActive] = useState(false);
    const [activeIdx, setActiveIdx] = useState<number>(-1);

    const siteName = "MimirLedger";
    return (
        <div>
            <NavTitle navLogo={siteName}/>
            <nav className={`nav`}>
                <div className={`${navActive ? "active" : ""} flex bg-slate-200`}>
                    {MENU_LIST.map((menu, idx) => (
                        <div className={`bg-blue-500 p-4 font-bold text-white text-xl`}
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
