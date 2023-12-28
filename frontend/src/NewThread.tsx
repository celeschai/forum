import { FormEvent, useState } from "react";
import { redirect } from "react-router-dom";
import { port } from "./App";

const NewThread = () => {
  const [title, setTitle] = useState('');
  const [tag, setTag] = useState('');
  const [author, setAuthor] = useState('mario');
  //const navigate = useNavigate();

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    const data = { 
      username: author, 
      title: title, 
      tag: tag, 
    };

    fetch("http://localhost:" + port + "/new/thread", {
      method: 'POST',
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data)
    }).then(() => {
      redirect("/feed");
    })
  }

  return (
    <div className="create">
      <h2>Add a New Blog</h2>
      <form onSubmit={handleSubmit}>
        <label>Thread title:</label>
        <input 
          type="text" 
          required 
          value={title}
          onChange={(e) => setTitle(e.target.value)}
        />
        <label>Your Username:</label>
        <textarea
          required
          value={author}
          onChange={(e) => setAuthor(e.target.value)}
        ></textarea>
        <label>Blog author:</label>
        <select
          value={tag}
          onChange={(e) => setTag(e.target.value)}
        >
          <option value="ut">University Town</option>
          <option value="soc">School of Computing</option>
        </select>
        <button>Add Blog</button>
      </form>
    </div>
  );
};

export default NewThread;