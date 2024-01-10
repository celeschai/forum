import { Link } from "react-router-dom";

const Navbar = () => {
    const title: string = "Foodie Gossips"

    return (
        <nav className="navbar">
            <div className="name">
                <h1>{ title }</h1>
            </div>
            <div className="links">
                <h2>
                <Link to="/feed/latest">Feed </Link>
                <Link to="/account">Account </Link>
                <Link to="/new/thread" >Create Thread</Link>
                </h2>
            </div>
        </nav>
    );
}
 
export default Navbar;
