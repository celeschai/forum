import { useState, useEffect } from 'react';

export type writingType = {
    id: number,
    title: string,
    content: string,
    username: string,
}

const useFetch = (url: string) => {
  const [data, setData] = useState<Array<writingType>|null>(null);
  const [isPending, setIsPending] = useState<boolean>(true);
  const [error, setError] = useState<any|null>(null);

  useEffect(() => {
    const abortCont = new AbortController();


    setTimeout(() => {
      fetch(url, { signal: abortCont.signal })
      .then(res => {
        if (!res.ok) { // error coming back from server
          throw Error('could not fetch the data for that resource');
        } 
        return res.json();
      })
      .then(data => {
        setIsPending(false);
        setData(data);
        setError(null);
      })
      .catch(err => {
        if (err.name === 'AbortError') {
          console.log('fetch aborted')
        } else {
          // auto catches network / connection error
          setIsPending(false);
          setError(err.message);
        }
      })
    }, 1000);

    // abort the fetch
    return () => abortCont.abort();
  }, [url])

  return { data, isPending, error };
}
 
export default useFetch;