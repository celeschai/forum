import { useState } from 'react';
import useFetch from '../Helpers/useFetch';
import { useParams } from 'react-router-dom';
import { handlePatch } from '../Helpers/handlers';

export const EditComment = ({ url }: { url: string }) => {
  const { id } = useParams()
  if (id == null) {
    throw Error("missing identity")
  }
  const { data, isPending, error } = useFetch(url.concat('/comment/', id))
  return (
    <div className="create">
      <h2>Edit Thread</h2>
      {isPending && <div>Loading...</div>}
      {error && <div>{error}</div>}
      {data && <EditCommentForm url={url} iniComment={data["content"]} id={id} />}
    </div>
  );
};

export default EditComment;


export const EditCommentForm = (
  { url, iniComment, id }: {
    url: string,
    iniComment: string,
    id: string,
  }) => {
  const [comment, setComment] = useState(iniComment);

  const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const req = { input1: comment, input2: null };

    handlePatch(url, 'comment', Number(id), req)
  }

  return (
    <form onSubmit={handleSubmit}>
      <label>Post title:</label>
      <textarea
        required
        value={comment}
        onChange={(e) => setComment(e.target.value)}
      />
      <button>Add Comment</button>
    </form>
  );
}

