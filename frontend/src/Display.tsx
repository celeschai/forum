import { Link } from 'react-router-dom';
import { useNavigate } from "react-router-dom";
import { Button } from '@mui/base';
import { JsonObjectExpression, JsonObjectExpressionStatement } from 'typescript';
 
export interface content {
    id: number,
    username: string,
    created: string,
}

export interface thread extends content {
    title: string,
    tag: string,
}

export interface post extends content {
    title: string,
    content: string,
}

// interface comment extends content {
//     content: string,
// }

export const DisplayThreads = ({list}: {list: Array<thread>}) => {
    return (
        <div className="list">
        {list.map(
            elem => 
                (<div className="preview" key={elem.id} >
                <Link to={`/thread/${elem.id}`}>
                    <h2>{ elem.title }</h2>
                    <h3>{ elem.tag }</h3>
                    <h4>by { elem.username }</h4>                  
                </Link>
                </div>
                )
        )}
        </div>
    );
}

export const DisplayPosts = ({list}: {list: Array<post>}) => {
    return (
        <div className="list">
        {list.map(
            elem => 
                (<div className="preview" key={elem.id} >
                <Link to={`/thread/${elem.id}`}>
                    <h2>{ elem.title }</h2>
                    <p>{ elem.content }</p>
                    <h3>created on {elem.created}</h3>
                    <h4>by { elem.username }</h4>                  
                </Link>
                </div>
                )
        )}
        </div>
    );
}
