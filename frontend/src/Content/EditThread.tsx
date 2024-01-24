import { useState } from 'react';
import useFetch from '../Helpers/useFetch';
import { useParams } from 'react-router-dom';
import { handlePatch } from '../Helpers/handlers';



export const EditThread = ({ url }: { url: string }) => {
    const { id } = useParams()
    if (id == null) {
        throw Error("missing identity")
    }
    const { data, isPending, error } = useFetch(url.concat('/user/thread/', id))
    return (
        <div className="create">
            <h2>Edit Thread</h2>
            {isPending && <div>Loading...</div>}
            {error && <div>{error}</div>}
            {data && <EditThreadForm url={url} iniTitle={data["title"]} iniTag={data["tag"]} id={id} />}
        </div>
    );
};

export default EditThread;


export const EditThreadForm = (
    { url, iniTitle, iniTag, id }: {
        url: string,
        iniTitle: string,
        iniTag: string,
        id: string,
    }) => {
    const [title, setTitle] = useState(iniTitle);
    const [tag, setTag] = useState(iniTag);

    const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        const req = { input1: title, input2: tag };

        handlePatch(url, '/thread', Number(id), req)
    }

    return (
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
                <option value="UT">University Town</option>
                <option value="SoC">School of Computing</option>
                <option value="FoS">Faculty of Science</option>
            </select>
            <button>Edit thread</button>
        </form>
    );
}

