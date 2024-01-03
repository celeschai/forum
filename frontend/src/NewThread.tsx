import { FormEvent, useState } from "react";
import { useNavigate } from "react-router-dom";

const NewThread = () => {
  const [title, setTitle] = useState('');
  const [tag, setTag] = useState('');

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    const data = { 
      title: title, 
      tag: tag, 
    };

    fetch("http://localhost:2999/newthread", {
      method: 'POST',
      headers: { 
        "Content-Type": "application/json",
        "Accept": "application/json", },
      credentials: 'include',
      body: JSON.stringify(data)
    }).then((response) =>  
        response.json()
        .then((stat) => {
          if (response.status === 200) {
            window.location.href = "/feed/latest";
          } else {
            console.log(stat.resp);
          }
      
    }))}
  

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
        {/* <label>Your Username:</label>
        <textarea
          required
          value={author}
          onChange={(e) => setAuthor(e.target.value)}
        ></textarea> */}
        <label>Location:</label>
        <select
          value={tag}
          onChange={(e) => setTag(e.target.value)}>
          <option value="University Town">University Town</option>
          <option value="School of Computing">School of Computing</option>
        </select>
        <button>Add Blog</button>
      </form>
    </div>
  );
};

export default NewThread;