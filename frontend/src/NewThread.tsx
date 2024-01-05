import { FormEvent, useState } from "react";
import { handlePost } from "./handlers";

const NewThread = ({url}:{url:string}) => {
  const [title, setTitle] = useState('');
  const [tag, setTag] = useState('');

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    const data = { 
      title: title, 
      tag: tag, 
    };

    handlePost(url, "/new/thread", data, "feed/latest")
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
        <label>Location:</label>
        <select
          value={tag}
          onChange={(e) => setTag(e.target.value)}>
          <option value="University Town">University Town</option>
          <option value="School of Computing">School of Computing</option>
        </select>
        <button>Add Thread</button>
      </form>
    </div>
  );
};

export default NewThread;