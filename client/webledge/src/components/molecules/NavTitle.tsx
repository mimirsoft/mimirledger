
export default function NavTitle( props:{navLogo:string})  {
    return (
        <div className="navlogo">
            <img className="m-2 navlogo-image font-bold text-4xl text-blue-500" src='/logo.png' alt={props.navLogo}></img>
        </div>
    )
}
