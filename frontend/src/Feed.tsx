import { useNavigate, useParams } from "react-router-dom";
import useFetch from "./useFetch";
import port from "./App";
import Display from "./Display";

const Feed = () => {
    //const { tag } = useParams();
    const url: string = 'http://localhost:' + port + '/feed';
    const { data, error, isPending } = useFetch(url);

    return (
        <div className="threads">
            <h2>Latest Threads</h2>
                { isPending && <div>Loading...</div> }
                { error && <div>{ error }</div> }
                { data && (
                    <article>
                       <Display list={data} />
                    </article>)
                }
        </div>
    );
}
 
export default Feed;