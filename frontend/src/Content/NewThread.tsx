import { FormEvent, useState } from "react";
import { handlePost } from "../Helpers/handlers";

const NewThread = ({ url }: { url: string }) => {
  const [title, setTitle] = useState('');
  const [tag, setTag] = useState('University Town');

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    const data = {
      title: title,
      tag: tag,
    };

    handlePost(url, "/new/thread", data, "/feed/latest")
  }


  return (
    <div className="createthread">
      <div className="create">
        <h2>Add a New Thread</h2>
        <form onSubmit={handleSubmit}>
          <div className="title">
            <label>Thread title:</label>
            <input
              type="text"
              required
              placeholder="(enter thread title)"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
            />
          </div>
          <div className="filter">
            <label>Location:</label>
            <select
              value={tag}
              onChange={(e) => setTag(e.target.value)}>
              <option value="University Town">University Town</option>
              <option value="School of Computing">School of Computing</option>
            </select>
          </div>
          <button>Add Thread</button>
        </form>
      </div>
    </div>

  );
};

export default NewThread;