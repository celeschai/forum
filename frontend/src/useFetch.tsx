import { useState, useEffect } from 'react';


const useFetch = (url: string) => {
  const [data, setData] = useState(null);
  const [isPending, setIsPending] = useState(true);
  const [error, setError] = useState(null);
  console.log(url);

  useEffect(() => {
    const abortCont = new AbortController();
    setTimeout(() => {
      fetch(url, 
        { 
          method: 'GET',
          headers: { "Content-Type": "application/json" },
          credentials: 'include',
          signal: abortCont.signal 
        })
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

    return () => abortCont.abort();
  }, [url])

  console.log(data);
  return { data, isPending, error };
}
 
export default useFetch;