import { FormEvent, useState } from "react";
import { handlePost } from "./handlers";

const NewComment = ({url, postid}: {url: string, postid: string}) => {
  const [comment, setComment] = useState('');
 
  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    const data = { 
      postid: postid,
      content: comment
    };

    handlePost(url, "/new/comment", data, "/postcomments/".concat(postid))
  }

  return (
    <div className="add">
      <h2>Add a New Comment</h2>
      <form onSubmit={handleSubmit}>
        <textarea 
          required 
          value={comment}
          onChange={(e) => setComment(e.target.value)}
        />
        <button>Add Comment</button>
      </form>
    </div>
  );
};

export default NewComment;