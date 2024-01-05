import { Link, Navigate } from 'react-router-dom';
import { handleDelete, handlePatch } from './handlers';
 
interface content {
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

export interface comment extends content {
    content: string,
}

export const DisplayThreads = ({
    url, list, allowEdit}: {
        url: string,
        list: Array<thread>;
        allowEdit: boolean;
}) => {


    return (    
        <div className="list">
        {list.map(
            (elem: thread) => 
                (<div className="thread" key={elem.id} >
                    <Link to={`/threadposts/${elem.id}`}>
                        <h2>{ elem.title }</h2>
                        <h3>{ elem.tag }</h3>
                        <h4>by { elem.username } on {elem.created}</h4>
                    </Link>
                    <h5>
                        { allowEdit && (
                            <div>
                                <button onClick={()=> handleDelete(url,"thread", elem.id)}>
                                    Delete
                                </button> 
                                <Link to = {("/thread/").concat(String(elem.id))}>
                                    <button> Edit </button> 
                                </Link>
                            </div>
                        )}
                    </h5>                  
                </div>)
        )}
        </div>
    )
};

export const DisplayPosts = ({
    url, list, allowEdit}: {
        url: string,
        list: Array<post>;
        allowEdit: boolean;
}) => {

    return (
        <div className="list">
        {list.map(
            (elem: post) => 
                (<div className="post" key={elem.id} >
                    <Link to={`/postcomments/${elem.id}`}>
                        <h2>{ elem.title }</h2>
                        <p>{ elem.content }</p>
                        <h3>created on { elem.created }</h3>
                        <h4>by { elem.username }</h4>                  
                    </Link>
                    <h5>
                        { allowEdit && (
                            <div>
                                <button onClick={()=> handleDelete(url,"thread", elem.id)}>
                                    Delete
                                </button> 
                                <Link to = {("/post/").concat(String(elem.id))}>
                                    <button> Edit </button> 
                                </Link>
                            </div>
                        )}
                    </h5>                  
                </div>)
        )}
        </div>
    );
}

export const DisplayComments = ({
    url, list, allowEdit}: {
        url: string,
        list: Array<comment>;
        allowEdit: boolean;
}) => {

    return(
        <div className="comments">
        {list.map(
            (elem: comment) => 
                (<div className="comment" key={elem.id} >
                        <p>{ elem.content }</p>
                        <h3>created on {elem.created}</h3>
                        <h4>by { elem.username }</h4>                  
                    <h5>
                        { allowEdit && (
                            <div>
                                <button onClick={()=> handleDelete(url,"thread", elem.id)}>
                                    Delete
                                </button> 
                                <Link to = {("/comment/").concat(String(elem.id))}>
                                    <button> Edit </button> 
                                </Link>
                            </div>
                        )}
                    </h5>    
                </div>
                )
        )}
        </div>
    );
}
