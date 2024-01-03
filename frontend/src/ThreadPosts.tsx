import { useParams } from "react-router-dom";
import useFetch from "./useFetch";
import { DisplayThreads, DisplayPosts, thread, post} from "./Display";

type threadposts = {
    thread: thread[],
    posts: post[],
}

const ThreadPosts = ({url}: {url: string}) => {

    const { id } = useParams();
    if (id === undefined) { 
        window.location.href = "/notfound";
    }
    const { data, error, isPending }: {
      data: threadposts | null,
      error: string | null,
      isPending: boolean,
    } = useFetch(String.prototype.concat(url, '/threadposts/', Number(id).toString()))

    return (
      <div className="Thread-Posts">
        <div className="Thread">
            { isPending && <div>Loading...</div> }
            { error && <div>{ error }</div> }
            { data && <DisplayThreads url={url} allowDel = {false} list={data['thread']} /> }
        </div>
        <div className="Post">
            { data && <DisplayPosts list={data['posts']} /> }
        </div>
    </div>
  );
}
 
export default ThreadPosts;