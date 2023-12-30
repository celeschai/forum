import useFetch from "./useFetch";
import { DisplayThreads } from "./Display";
import { useParams } from "react-router-dom";

//figure out how to make pages

const Feed = () => {
    
    const { tag } = useParams<{tag: string}>();
  
    const { data, error, isPending } = useFetch('http://localhost:2999/feed/' + tag);

    return (
        <div className="threads">
            <h1>Latest Threads</h1>
                { isPending && <div>Loading...</div> }
                { error && <div>{ error }</div> }
                { data && (
                    <article>
                       <DisplayThreads list={data} />
                    </article>)}
        </div>
    );
}
 
export default Feed;