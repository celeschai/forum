import useFetch from "../Helpers/useFetch";
import { DisplayThreads, DisplayPosts, DisplayComments, thread, post, comment } from "../Helpers/Display";
import { handlePost } from "../Helpers/handlers"

type accountPosts = {
    username: string,
    threads: thread[],
    posts: post[],
    comments: comment[],
}

const Account = ({ url }: { url: string }) => {
    const { data, error, isPending }: {
        data: accountPosts | null,
        error: string | null,
        isPending: boolean,
    } = useFetch(url.concat('/account'));

    return (
        <div className="account">
            {isPending && <div>Loading...</div>}
            {error && <div>{error}</div>}
            {!data && <button onClick={() => (window.location.href = "/login")}>Sign in</button>}
            {data && <User url={url} data={data} />}
        </div>
    );
}

export default Account;

const User = ({ url, data }: { url: string, data: accountPosts }) => {
    const handleSignOut = () => {
        if (window.confirm("Are you sure you want to log out of this account? You will have to sign in again.")) {
            handlePost(url, "/signout", null, "")
        }
    }
    return (
        <div className="user">
            <h3>You are logged in as: {data["username"]}</h3>
            <button onClick={handleSignOut}>Sign out</button>

            <div className="Thread">
                <h1>Threads created</h1>
                <DisplayThreads url={url} allowEdit={true} list={data["threads"]} />
            </div>

            <div className="Post">
                <h1>Posts created</h1>
                <DisplayPosts url={url} allowEdit={true} list={data["posts"]} />
            </div>

            <div className="Comment">
                <h1>Comments created</h1>
                <DisplayComments url={url} allowEdit={true} list={data["comments"]} />
            </div>
        </div>
    );
}

