import { useParams } from "react-router-dom";
import useFetch from "./useFetch";
import { DisplayThreads, DisplayPosts, DisplayComments, thread, post, comment} from "./Display";
import { useState } from "react";
import NewPost from "./NewPost";
import NewComment from "./NewComment";

type threadposts = {
    parent: thread,
    child: post[],
}

type postcomments = {
    root: thread,
    parent: post,
    child: comment[],
}

export const ThreadPosts = ({url}: {url: string}) => {
    const [add, setAdd] = useState(false);
    const { id } = useParams();
    if (id === undefined) { 
       throw Error("thread not found")
    }
    const { data, error, isPending }: {
      data: threadposts | null,
      error: string | null,
      isPending: boolean,
    } = useFetch(String.prototype.concat(url, '/parentchild/thread/', id))

    return (
      <div className="Thread-Posts">
        <div className="Thread">
            { isPending && <div>Loading...</div> }
            { error && <div>{ error }</div> }
            { data && <DisplayThreads url={url} allowEdit = {false} list={data['parent']} /> }
            { !add && <button onClick={() => setAdd(true)}>Add Post</button> }
        </div>
        <br></br>
        { add && <NewPost url={url} threadid={id}/>}
        <br></br>
        <div className="Post">
            { data && <DisplayPosts url={url} allowEdit = {false} list={data['child']} /> }
        </div>
    </div>
  );
}
 
export const PostComments = ({url}: {url: string}) => {
  const [add, setAdd] = useState(false);
  const { id } = useParams();
  if (id === undefined) { 
      throw Error("post not found")
  }
  const { data, error, isPending }: {
    data: postcomments | null,
    error: string | null,
    isPending: boolean,
  } = useFetch(String.prototype.concat(url, '/parentchild/post/', Number(id).toString()))

  return (
    <div className="Post-Comments">
      <div className="Thread">
          { isPending && <div>Loading...</div> }
          { error && <div>{ error }</div> }
          { data && <DisplayThreads url={url} allowEdit = {false} list={data['root']} /> }
      </div>
      <br></br>
      <div className="Post">
          { isPending && <div>Loading...</div> }
          { error && <div>{ error }</div> }
          { data && <DisplayPosts url={url} allowEdit = {false} list={data['parent']} /> }
          { !add && <button onClick={() => setAdd(true)}>Add Post</button> }
      </div>
      <br></br>
      { add && <NewComment url={url} postid={id}/>}
      <br></br>
      <div className="Comment">
          { data && <DisplayComments url={url} allowEdit = {false} list={data['child']} /> }
      </div>
  </div>
);
}
