import useFetch from "../Helpers/useFetch";
import { DisplayThreads } from "../Helpers/Display";
import { useState } from "react";

const Feed = ({ url }: { url: string }) => {

    const [tag, setTag] = useState('latest');

    const { data, error, isPending } = useFetch(String.prototype.concat(url, '/feed/', tag));

    return (
        <div className="feed">
            <div className="filter">
                <h1>Threads filtered by</h1>
                <select value={tag}
                    onChange={(e) => setTag(e.target.value)}>
                    <option value="latest">Latest</option>
                    <option value="UT">University Town</option>
                    <option value="SoC">School of Computing</option>
                    <option value="FoS">Faculty of Science</option>
                    
                </select>
            </div>
            {isPending && <div>Loading...</div>}
            {error && <div>{error}</div>}
            {data && (
                <article>
                    <DisplayThreads url={url} list={data} allowEdit={false} />
                </article>)}
        </div>
    );
}

export default Feed;