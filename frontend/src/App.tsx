import './App.css';
import SignIn from './SignIn';
import SignUp from './SignUp';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import Navbar from './Navbar';
import NotFound from './NotFound';
import Feed from './Feed';
import NewThread from './NewThread';
import Account from './Account';
import { ThreadPosts, PostComments } from './ParentChild';
import Start from './Start';
import { EditThread } from './EditThread';
import { EditPost } from './EditPost';
import { EditComment } from './EditComment';



export const hosturl: string = "http://localhost:2999";

function App() {

  return (
      <Router>
        <div className="app">
          <Navbar />

          <Routes>
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
