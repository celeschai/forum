import { Link } from "react-router-dom";

const Navbar = () => {
    const title: string = "Foodie Gossips"

    return (
        <nav className="navbar">
            <h1>{ title }</h1>
            <div className="links">
                <Link to="/feed/latest">Latest</Link>
                <Link to="/new/thread" style={{ 
                color: 'white', 
                backgroundColor: '#f1356d',
                borderRadius: '8px' 
                }}>Create Thread</Link>
                <Link to="/account">Account</Link>
            </div>
        </nav>
    );
}
 
export default Navbar;