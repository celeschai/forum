import {thread, post, comment} from './Display'
import { useState, FormEvent } from 'react';
import useFetch from './useFetch';
import { useParams } from 'react-router-dom';

export const handleDelete = (url: string, type: string, id: number) => {
    //console.log(String.prototype.concat(url, '/delete/', type, '/', id.toString())
    if (window.confirm("Are you sure you want to delete this?")) {
        fetch(
            String.prototype.concat(url, '/delete/', type, '/', id.toString()), { 
            method: 'DELETE',
            headers: { 
                "Content-Type": "application/json",
                "Accept": "application/json", },
            credentials: 'include',
        }).then((resp) => {
            if (resp.ok) {
                window.location.href = "/account"  
            } else if (resp.status === 401) {
                console.log ("Log in this account to perform this action")
                window.location.href = "/login"
            } else {
                console.log ("error performing action")
                window.location.reload()
            }
        })
    } else {
        console.log("Request cancelled")
    }
}

type patchRequest = {
    input1: string | FormDataEntryValue | null,
    input2: string | FormDataEntryValue | null,
}

export const handlePatch = (url: string, type: string, id: number, data: patchRequest) => {
    //console.log(String.prototype.concat(url, '/delete/', type, '/', id.toString())
    if (window.confirm("Are you sure you want to delete this?")) {
        fetch(
            String.prototype.concat(url, '/patch/', type, '/', id.toString()), { 
            method: 'PATCH',
            headers: { 
                "Content-Type": "application/json",
                "Accept": "application/json", },
            credentials: 'include',
            body: JSON.stringify(data)
        }).then((resp) => {
            if (resp.ok) {
                window.location.href = "/account"  
            } else if (resp.status === 401) {
                console.log ("Log in this account to perform this action")
                window.location.href = "/login"
            } else {
                console.log ("error performing action")
                window.location.reload()
            }
        })
    } else {
        console.log("Request cancelled")
    }
}

export const EditThread = ({url}: {url:string}) => {
    const [title, setTitle] = useState('');
    const [tag, setTag] = useState('');

    const {id} = useParams()
    if (id == null) {
        window.location.href = "/notfound"
    }

    const {data, isPending, error} = useFetch(url.concat('/thread/', String(id)))
    if (error == null && isPending && (data != null)) {
        setTitle (data["title"])
        setTag (data["tag"])
    } 
    else if (error != null) {
        throw error 
    }

    const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        const form = new FormData(event.currentTarget);
        const req = { input1: form.get('title'), input2: form.get('tag') };

        handlePatch(url, 'thread', Number(id), req)
    }
  
    return (
      <div className="create">
        <h2>Edit Thread</h2>
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
          <button>Add Blog</button>
        </form>
      </div>
    );
  };
  
  export default EditThread;