import './App.css';
import SignIn from './SignIn';
import SignUp from './SignUp';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import Navbar from './Navbar';
import NotFound from './NotFound';
import Feed from './Feed';
import NewThread from './NewThread';
import Account from './Account';
import NewPost from './NewPost';
import NewComment from './NewComment';
import ThreadPosts from './ThreadPosts';
import Start from './Start';
import {EditThread} from './Edit';
import {thread} from './Display';

export const hosturl: string = "http://localhost:2999"; 

function App() {
  
  return (
  <Router>
      <div className="App">
        <Navbar />
        <div className="content">
          <Routes>
              <Route path="/home" element = {<Start url={hosturl.concat("/home")}/>} />

              <Route path="/login" element = {<SignIn url={hosturl.concat("/login")} />}/>

              <Route path="/feed/:tag" element = {<Feed url={hosturl} />}/>

              <Route path="/signup" element = {<SignUp url={hosturl.concat("/signup")}/>} />

              <Route path="/account" element = {<Account url={hosturl} />} />

              <Route path="/thread/:id" element = {<EditThread url = {hosturl}/>} />

              <Route path="/threadposts/:id" element = {<ThreadPosts url={hosturl} />} />

              <Route path="/new/thread" element = {<NewThread url={hosturl} />} />

              <Route path="/new/post" element = {<NewPost />} />

              <Route path="/new/comment" element = {<NewComment />} />

              <Route path="*" element = {<NotFound />} />
              
          </Routes>
        </div>
      </div>
    </Router>
  );
}

export default App;
