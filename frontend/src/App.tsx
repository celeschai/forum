import SignIn from './User/SignIn';
import SignUp from './User/SignUp';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import Navbar from './User/Navbar';
import NotFound from './Content/NotFound';
import Feed from './User/Feed';
import NewThread from './Content/NewThread';
import Account from './User/Account';
import { ThreadPosts, PostComments } from './Content/ParentChild';
import Start from './User/Start';
import { EditThread } from './Content/EditThread';
import { EditPost } from './Content/EditPost';
import { EditComment } from './Content/EditComment';

const HTTP: string = process.env.REACT_APP_HTTPS === "true" ? "https" : "http";
const host: string = process.env.REACT_APP_HOST != null ? process.env.REACT_APP_HOST : "localhost"
const serverport: string = process.env.REACT_APP_BACK_PORT != null ? process.env.REACT_APP_BACK_PORT : "3000"

const hosturl: string = HTTP.concat("://", host, ":", serverport)


function App() {

  return (
    <Router>
      <div className="app">
        <Navbar />

        <Routes> //urls can also be imported, passing in as props so components are more resuable
          <Route path="" element={<Start url={hosturl.concat("/home")} />} />

          <Route path="/login" element={<SignIn url={hosturl.concat("/login")} initialResult={'Login required'} />} />

          <Route path="/feed/:tag" element={<Feed url={hosturl} />} />

          <Route path="/signup" element={<SignUp url={hosturl.concat("/signup")} initialResult={''} />} />

          <Route path="/account" element={<Account url={hosturl} />} />

          <Route path="/thread/:id" element={<EditThread url={hosturl} />} />

          <Route path="/post/:id" element={<EditPost url={hosturl} />} />

          <Route path="/comment/:id" element={<EditComment url={hosturl} />} />

          <Route path="/threadposts/:id" element={<ThreadPosts url={hosturl} />} />

          <Route path="/postcomments/:id" element={<PostComments url={hosturl} />} />

          <Route path="/new/thread" element={<NewThread url={hosturl} />} />

          <Route path="*" element={<NotFound />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;
