import { useParams } from "react-router-dom";

const Account = () => {
    const { username } = useParams();
    


    return (
        <div className="account"> 
            <h2>Profile of {username}</h2>
            <h3>Threads</h3>

            <h3>Posts</h3>

            <h3>Comments</h3>
        </div>
    );
}
 
export default Account;