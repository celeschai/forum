import useFetch from "./useFetch";
import { DisplayThreads } from "./Display";
import { useState } from "react";
//figure out how to make pages

const Feed = ({url}: {url: string}) => {

    const [tag, setTag] = useState('latest');
  
    const { data, error, isPending } = useFetch(String.prototype.concat(url, '/feed/', tag));

    return (
        <div className="threads">
            <h1>Threads filtered by</h1>
                <select value={tag}
                    onChange={(e) => setTag(e.target.value)}>
                    <option value="latest">Latest</option>
                    <option value="University Town">University Town</option>
                    <option value="School of Computing">School of Computing</option>
                </select>
                { isPending && <div>Loading...</div> }
                { error && <div>{ error }</div> }
                { data && (
                    <article>
                       <DisplayThreads url = {url} list = {data} allowEdit = {false} />
                    </article>)}
        </div>
    );
}
 
export default Feed;