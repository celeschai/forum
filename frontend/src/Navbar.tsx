import { Link } from "react-router-dom";

const Navbar = () => {
    const title: string = "Foodie Gossips"

    return (
        <nav className="navbar">
            <h1>{ title }</h1>
            <div className="links">
                <Link to="/feed/latest">Feed</Link>
                <Link to="/account">Account</Link>
                <Link to="/newthread" style={{ 
                color: 'white', 
                backgroundColor: '#f1356d',
                borderRadius: '8px' 
                }}>Create Thread</Link>
            </div>
        </nav>
    );
}
 
export default Navbar;