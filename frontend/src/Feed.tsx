import useFetch from "./useFetch";
import hosturl from "./App";
import Display from "./Display";

type writingType = {
    id: number,
    title: string,
    username: string,
    content: string,
    created_at: Date,
}


const Feed = () => {
    const url: string = hosturl + '/feed/latest/0';
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