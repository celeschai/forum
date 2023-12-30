import { FormEvent, useState } from "react";
import { useNavigate } from "react-router-dom";

const NewThread = () => {
  const [title, setTitle] = useState('');
  const [tag, setTag] = useState('');
  const [author, setAuthor] = useState('mario');
  const navigate = useNavigate();

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    const data = { 
      title: title, 
      username: author, 
      tag: tag, 
    };

    fetch("http:localhost:2999/new/thread", {
      method: 'POST',
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data)
    }).then(() => {
      navigate("/feed");
    })
  }

  return (
    <div className="create">
      <h2>Add a New Thread</h2>
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
        <label>Location:</label>
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