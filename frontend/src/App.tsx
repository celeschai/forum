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

export const port: string = "2999"; 

function App() {
  
  return (
  <Router>
      <div className="App">
        <Navbar />
        <div className="content">
          <Routes>
              <Route path="/signin" element = {<SignIn />}/>

              <Route path="/feed" element = {<Feed />}/>

              <Route path="/newaccount" element = {<SignUp />} />

              <Route path="/account/:username" element = {<Account />} />

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
