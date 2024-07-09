import React from 'react'
import NavBar from "../organisms/NavBar";
import MainContainer from "../organisms/MainContainer";
import Footer from "../organisms/Footer";
import { Outlet } from "react-router-dom";

export default function OuterContainer() {
    return (
        <><div className="w-full">
            <NavBar />
            <MainContainer><Outlet /></MainContainer>
            <Footer/>
        </div>
        </>
    )
}
