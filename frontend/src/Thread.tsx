import { useParams } from "react-router-dom";
import useFetch from "./useFetch";
import {DisplayPosts, thread, post} from "./Display";

type threadposts = {
    thread: thread[],
    posts: post[],
}

const Thread = ({allowdelete}: {allowdelete: boolean}) => {

    const { id } = useParams();
    const { data, error, isPending }: {
      data: threadposts | null,
      error: string | null,
      isPending: boolean,
    } = useFetch('http://localhost:2999/threadposts/' + id);
  

    const handleClick = 
      () => { if (data !== null) 
          {
              fetch('http://localhost:2999' + data['id'], {
                  method: 'DELETE'
              }).then(() => {
                  window.location.href = "/feed/account/" + data['username'] 
              }) 
          }
      }

    return (
      <div className="Thread-Posts">
        <div className="Thread">
            { isPending && <div>Loading...</div> }
            { error && <div>{ error }</div> }
            { data && <DisplayPosts list={data['thread']} /> }
            { allowdelete && <button onClick={handleClick}>Delete Thread</button> }
        </div>
        <div className="Post">
            { data && <DisplayPosts list={data['posts']} /> }
        </div>
    </div>
  );
}
 
export default Thread;