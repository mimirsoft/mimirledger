import {Link } from "react-router-dom";

const NavButton = ( props:{href:string, text:string, active:boolean} ) => {
    return (
        <Link to={props.href} className={`nav__item ${props.active ? "active" : ""}`}>
            {props.text}
        </Link>
    );
};

export default NavButton;