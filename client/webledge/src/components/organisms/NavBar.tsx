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
                <div></div>
                <div
                    onClick={() => setNavActive(!navActive)}
                    className={`nav__menu-bar`}
                >
                    <div></div>
                    <div></div>
                    <div></div>
                </div>
                <div className={`${navActive ? "active" : ""} nav__menu-list`}>
                    {MENU_LIST.map((menu, idx) => (
                        <div className={`nav__menu-listbox-item`}
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
