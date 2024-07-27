import {Link, useMatch, useResolvedPath} from "react-router-dom";

const NavButton = ( props:{href:string, text:string} ) => {
    const resolvedPath = useResolvedPath(props.href)
    const isActive = useMatch({path:resolvedPath.pathname, end:true})
    return (
        <Link to={props.href} className={`${isActive ? "underline" : ""} p-4 hover:underline`}>
            {props.text}
        </Link>
    );
};

export default NavButton;