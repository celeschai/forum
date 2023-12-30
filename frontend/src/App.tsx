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
import Thread from './Thread';

export const hosturl: string = "http://localhost:2999"; 

function App() {
  
  return (
  <Router>
      <div className="App">
        <Navbar />
        <div className="content">
          <Routes>
              <Route path="/signin" element = {<SignIn />}/>

              <Route path="/feed/:tag" element = {<Feed />}/>

              <Route path="/newaccount" element = {<SignUp />} />

              <Route path="/account/:username" element = {<Account />} />

              <Route path="/thread/:id" element = {<Thread allowdelete = {false} />} />

              <Route path="/new/thread" element = {<NewThread />} />

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
