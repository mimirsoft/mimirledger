import React from 'react'

export default function NavTitle( props:{navLogo:string})  {
    return (
        <div className="navlogo">
            <img className="navlogo-image" src='/logo.png' alt={props.navLogo}></img>
        </div>
    )
}
