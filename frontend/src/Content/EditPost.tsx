import { useState } from 'react';
import useFetch from '../Helpers/useFetch';
import { useParams } from 'react-router-dom';
import { handlePatch } from '../Helpers/handlers';



export const EditPost = ({ url }: { url: string }) => {
    const { id } = useParams()
    if (id == null) {
        throw Error("missing identity")
    }
    const { data, isPending, error } = useFetch(url.concat('/post/', id))
    return (
        <div className="create">
            <h2>Edit Thread</h2>
            {isPending && <div>Loading...</div>}
            {error && <div>{error}</div>}
            {data && <EditCommentForm url={url} iniTitle={data["title"]} iniContent={data["content"]} id={id} />}
        </div>
    );
};

export default EditPost;


export const EditCommentForm = (
    { url, iniTitle, iniContent, id }: {
        url: string,
        iniTitle: string,
        iniContent: string,
        id: string,
    }) => {
    const [title, setTitle] = useState(iniTitle);
    const [content, setContent] = useState(iniContent);

    const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        const req = { input1: title, input2: content };

        handlePatch(url, 'post', Number(id), req)
    }

    return (
        <form onSubmit={handleSubmit}>
            <label>Post title:</label>
            <input
                type="text"
                required
                value={title}
                onChange={(e) => setTitle(e.target.value)}
            />
            <label>Body:</label>
            <textarea
                required
                value={content}
                onChange={(e) => setContent(e.target.value)} />
            <button>Add Post</button>
        </form>
    );
}

