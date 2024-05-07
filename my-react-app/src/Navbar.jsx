import { Link } from "react-router-dom"

function Navbar(props) {
    if (!props.isLogin) {
        return (
    <nav className='navbar navbar-text bg-body-tertiary justify-content-center'>
        <Link className='nav-link me-2' to="/component/login">Login</Link>
        <Link className='nav-link' to="/register">Register</Link>
    </nav>)
    }
    return (
        <></>
    )
}

export default Navbar