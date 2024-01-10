import { useParams } from "react-router-dom";
import useFetch from "./useFetch";
import { DisplayThreads, DisplayPosts, DisplayComments, thread, post, comment } from "./Display";
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

export const ThreadPosts = ({ url }: { url: string }) => {
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
    <div className="parentchild">
      <div className="parent_pc">
        {isPending && <div>Loading...</div>}
        {error && <div>{error}</div>}
        {data && <DisplayThreads url={url} allowEdit={true} list={data['parent']} />}

        <br></br>

        {!add && <button onClick={() => setAdd(true)}>Add Post</button>}
        {add && <NewPost url={url} threadid={id} />}
      </div>

      <div className="child_pc">
        {data && <DisplayPosts url={url} allowEdit={false} list={data['child']} />}
      </div>
    </div>
  );
}

export const PostComments = ({ url }: { url: string }) => {
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
    <div className="parentchild">
      <div className="parent_pc">
        {isPending && <div>Loading...</div>}
        {error && <div>{error}</div>}
        {data && <DisplayThreads url={url} allowEdit={true} list={data['root']} />}

        ↑

        {isPending && <div>Loading...</div>}
        {error && <div>{error}</div>}
        {data && <DisplayPosts url={url} allowEdit={true} list={data['parent']} />}

        <br></br>

        {!add && <button onClick={() => setAdd(true)}>Add Comment</button>}
        {add && <NewComment url={url} postid={id} />}
      </div>

      <div className="child_pc">
        {data && <DisplayComments url={url} allowEdit={true} list={data['child']} />}
      </div>
    </div>
  );
}
