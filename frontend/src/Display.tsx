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


function timeConverter(datestring: string) {
    const isoDateString: string = datestring;
    const isoDate = new Date(isoDateString);
    const year = isoDate.getUTCFullYear();
    const month = isoDate.getUTCMonth() + 1; // Months are zero-based, so we add 1
    const day = isoDate.getUTCDate();

    const formattedDate = `${year}-${month < 10 ? '0' : ''}${month}-${day < 10 ? '0' : ''}${day}`;
    return formattedDate;
}

export const DisplayThreads = ({
    url, list, allowEdit }: {
        url: string,
        list: Array<thread>;
        allowEdit: boolean;
    }) => {


    return (
        <div className="list">
            {list.map(
                (elem: thread) =>
                    <div className="parent">
                        <div className="content" key={elem.id} >
                            <Link to={`/threadposts/${elem.id}`}>
                                <h2>{elem.title}</h2>
                                <h3> 📍  {elem.tag}</h3>
                                <h4>by {elem.username} on {timeConverter(elem.created)}</h4>
                            </Link>
                            <h5>
                                {allowEdit && (
                                    <div className="buttons">
                                        <button onClick={() => handleDelete(url, "/thread", elem.id)}>
                                            🗑️
                                        </button>
                                        <Link to={("/thread/").concat(String(elem.id))}>
                                            <button> ✏️ </button>
                                        </Link>
                                    </div>
                                )}
                            </h5>
                        </div>
                    </div>
            )}
        </div>
    )
};

export const DisplayPosts = ({
    url, list, allowEdit }: {
        url: string,
        list: Array<post>;
        allowEdit: boolean;
    }) => {

    return (
        <div className="list">
            {list.map(
                (elem: post) =>
                    <div className="parent">
                        <div className="content" key={elem.id} >
                            <Link to={`/postcomments/${elem.id}`}>
                                <h3>{elem.title}</h3>
                                <p>{elem.content}</p>
                                <h4>by {elem.username} on {timeConverter(elem.created)}</h4>
                            </Link>
                            <h5>
                                {allowEdit && (
                                    <div>
                                        <button onClick={() => handleDelete(url, "thread", elem.id)}>
                                            Delete
                                        </button>
                                        <Link to={("/post/").concat(String(elem.id))}>
                                            <button> Edit </button>
                                        </Link>
                                    </div>
                                )}
                            </h5>
                        </div>
                    </div>
            )}
        </div>
    );
}

export const DisplayComments = ({
    url, list, allowEdit }: {
        url: string,
        list: Array<comment>;
        allowEdit: boolean;
    }) => {

    return (
        <div className="list">
            {list.map(
                (elem: comment) =>
                (<div className="content" key={elem.id} >
                    <p>{elem.content}</p>
                    <h4>by {elem.username} on {timeConverter(elem.created)}</h4>
                    <h5>
                        {allowEdit && (
                            <div>
                                <button onClick={() => handleDelete(url, "thread", elem.id)}>
                                    Delete
                                </button>
                                <Link to={("/comment/").concat(String(elem.id))}>
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
