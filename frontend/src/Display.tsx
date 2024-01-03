import { ClassNames } from '@emotion/react';
import { Link } from 'react-router-dom';
import { useParams } from 'react-router-dom';
import { Type } from 'typescript';

 
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

export interface comment extends content {
    content: string,
}

export const handleDelete = (url: string, type: string, id: number) => {
    //console.log(String.prototype.concat(url, '/delete/', type, '/', id.toString())
    if (window.confirm("Are you sure you want to delete this?")) {
        fetch(
            String.prototype.concat(url, '/delete/', type, '/', id.toString()), { 
            method: 'DELETE',
            headers: { 
                "Content-Type": "application/json",
                "Accept": "application/json", },
            credentials: 'include',
        }).then((resp) => {
            if (resp.ok) {
                window.location.href = "/account"  
            } else if (resp.status === 401) {
                console.log ("Log in this account to perform this action")
                window.location.href = "/login"
            } else {
                console.log ("Error deleting")
                window.location.reload()
            }
        })
    } else {
        console.log("Delete cancelled")
    }
}

export const DisplayThreads = ({
    url, list, allowDel}: {
        url: string,
        list: Array<thread>;
        allowDel: boolean;
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
                    { allowDel && ( <button onClick={()=> handleDelete(url,"thread", elem.id)}>Delete</button> ) }
                </h5>                  
                </div>
                )
        )}
        </div>
    )
};

export const DisplayPosts = ({list}: {list: Array<post>}) => {
    return (
        <div className="list">
        {list.map(
            (elem: post) => 
                (<div className="preview" key={elem.id} >
                <Link to={`/post/${elem.id}`}>
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

export const DisplayComments = ({list}: {list: Array<comment>}) => {
    return (
        <div className="comments">
        {list.map(
            (elem: comment) => 
                (<div className="preview" key={elem.id} >
                <Link to={`/comment/${elem.id}`}>
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
