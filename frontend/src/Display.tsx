import { Link } from 'react-router-dom';
import { useNavigate } from "react-router-dom";
import { Button } from '@mui/base';
 
// type displayProps = {
//     id: number,
//     title: string,
//     content: string,
//     username: string,
// }

const Display = ({list}: {list: Array<>}) => {
    const navigate = useNavigate();

    return (
        <div className="blog-list">
        {list.map(elem => (
            <div className="preview" key={elem.id} >
            <Link to={`/feed/${elem.id}`}>
                <h2>{ elem.title }</h2>
                <h3>{ elem.content }</h3>
                <p>by { elem.username }</p>                  
            </Link>
            </div>
        ))}
        </div>
    );
}
 
export default Display;