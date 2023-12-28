import { useNavigate } from "react-router-dom";

const Delete = (url: string) => {
    type deleteFunc = (url: string) => void;
    const navigate = useNavigate();

    const handleDelete: deleteFunc = 
    (url: string) => {
        fetch(url, {method: 'DELETE'})
        .then(() => {
            navigate('/');
            console.log('deleted');}) 
    }

    return ( 
        <button onClick={() => handleDelete(url)}>
            delete
        </button>     
     );
}
 
export default Delete;