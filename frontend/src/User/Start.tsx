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
            <h2>Foodie Gossips is a forum for foodies to share their thoughts on campus dining.</h2>
            <br></br>
            <Button onClick={handleclick}>Enter</Button>
            <br></br>
            <br></br>
            <p>
                This page is works best on a computer using Google Chrome.
            </p>
            <br></br>
            <p>
                Please enable Third Party Cookies and include this site * https://forum-front-ynvw.onrender.com * in your browser's cookies allowed site list.
            </p>
            <p>
                On mobile devices, please disable block all cookies and allow cross site tracking.
            </p>
            <a href = "https://github.com/celeschai/forum/blob/main/UserManual.pdf">user manual</a>
        </div>
    );
}

export default Start;