import { useParams } from "react-router-dom";
import useFetch from "./useFetch";
import {DisplayThreads, DisplayPosts, DisplayComments, thread, post, comment} from "./Display";

// type account = {
//     username: string,
//     email: string,
//}
type accountPosts = {
    username: string,
    threads: thread[],
    posts: post[],
    comments: comment[],
}

const Account = ({url}: {url: string}) => {
    // const handleClick = 
    //   () => { if (data !== null) 
    //       {
    //           fetch('http://localhost:2999' + data['id'], {
    //               method: 'DELETE'
    //           }).then(() => {
    //               window.location.href = "/feed/account/" + data['username'] 
    //           }) 
    //       }
    //   }

    const {data, error, isPending}: {
        data: accountPosts | null,
        error: string | null,
        isPending: boolean,
    } = useFetch(url.concat('/account'));


    return (
        <div className="account"> 
            <h1>Profile of {data && data["username"]}</h1>
            <h2>Threads created</h2>
                { isPending && <div>Loading...</div> }
                { error && <div>{ error }</div> }
                { data && <DisplayThreads url={url} allowEdit = {true} list = {data["threads"]} /> }

            <h2>Posts created</h2>
                { isPending && <div>Loading...</div> }
                { error && <div>{ error }</div> }
                { data && <DisplayPosts list = {data["posts"]} /> }

            <h2>Comments created</h2>
                { isPending && <div>Loading...</div> }
                { error && <div>{ error }</div> }
                { data && <DisplayComments list = {data["comments"]} /> }
        </div>
    );
}
 
export default Account;