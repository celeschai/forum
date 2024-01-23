import { Button } from "@mui/material";
import { useEffect } from "react";
import { useState } from "react";

const Start = ({ url }: { url: string }) => {
    const [check, setCheck] = useState(false);
    useEffect(() => {
        fetch(url,
            {
                method: 'GET',
                headers: { "Content-Type": "application/json" },
                credentials: 'include',
            })
            .then(res => {
                if (res.status === 200) { // error coming back from server
                    window.location.href = "/feed/latest";
                } else if (res.status === 401) {
                    window.location.href = "/login";
                } else {
                    throw Error('Something went wrong');
                }
            })
    }, [check])


    return (
        <div className="start">
            <h1>Welcome to Foodie Gossips!</h1>
            <p>Foodie Gossips is a forum for foodies to share their thoughts on campus dining.</p>
            <br></br>
            <Button onClick={() => setCheck(true)}>Enter</Button>
        </div>
    );
}

export default Start;