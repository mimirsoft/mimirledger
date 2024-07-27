
type MainContainerProps = {
    children:  React.ReactNode
}
export default function MainContainer({children}: MainContainerProps) {
    return (
        <div className="main-container">{children}
        </div>
    )
}
