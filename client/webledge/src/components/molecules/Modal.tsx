import React from 'react'

const  Modal = (props:{ open:boolean})  => {
    if ( !props.open)  return null

    return (
        <div>
           Modal
        </div>
    );
};

export default Modal;