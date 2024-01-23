import { Button } from "@mui/material";

const Start = ({ url }: { url: string }) => {
    const handleclick = () => {
        fetch(url,
            {   method: 'GET',
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
    }


    return (
        <div className="start">
            <h1>Welcome to Foodie Gossips!</h1>
            <p>Foodie Gossips is a forum for foodies to share their thoughts on campus dining.</p>
            <br></br>
            <Button onClick={handleclick}>Enter</Button>
        </div>
    );
}

export default Start;