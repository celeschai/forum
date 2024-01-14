import { FormEvent, useState } from "react";
import { handlePost } from "../Helpers/handlers";

const NewPost = ({ url, threadid }: { url: string, threadid: string }) => {
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    const data = {
      threadid: threadid,
      title: title,
      content: content,
    };

    handlePost(url, "/new/post", data, "/threadposts/".concat(threadid))
  }

  return (
    <div className="create">
      <h2>Add a New Post</h2>
      <form onSubmit={handleSubmit}>
        <label>Post title:</label>
        <input
          type="text"
          required
          value={title}
          onChange={(e) => setTitle(e.target.value)}
        />
        <label>Body</label>
        <textarea
          required
          value={content}
          onChange={(e) => setContent(e.target.value)} />
        <button>Add Post</button>
      </form>
    </div>
  );
};

export default NewPost;