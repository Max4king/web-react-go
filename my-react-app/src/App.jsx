import React from 'react';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import LoginPage from './LoginPage';
import './App.css'
import RegisterPage from './RegisterPage';
import { useState } from 'react';
// import Dashboard from './Dashboard';

function App() {
  const [isLogin, setIsLogin] = useState(false);
  
  if (isLogin) {
    return (
    <div>
    <nav className='navbar navbar-text bg-body-tertiary justify-content-center'>
        <Link className='nav-link me-2' to="/dashboard">Dashboard</Link>
    </nav>
    <Routes>
        <Route path="/" element={<LoginPage isLogin={isLogin} setIsLogin={setIsLogin}/>} />
        <Route path="/register" element={<RegisterPage isLogin={isLogin} setIsLogin={setIsLogin}/>} />
        {/* <Route path="/dashboard" element={<Dashboard isLogin={isLogin} setIsLogin={setIsLogin}/>} /> */}
      </Routes>
    </div>
    )
  }

  return (
    <div>
    <nav className='navbar navbar-text bg-body-tertiary justify-content-center'>
        <Link className='nav-link me-2' to="/dashboard">Dashboard</Link>
        <Link className='nav-link me-2' to="/">Login</Link>
        <Link className='nav-link' to="/register">Register</Link>
    </nav>
    <Routes>
        <Route path="/" element={<LoginPage isLogin={isLogin} setIsLogin={setIsLogin}/>} />
        <Route path="/register" element={<RegisterPage isLogin={isLogin} setIsLogin={setIsLogin}/>} />
        {/* <Route path="/dashboard" element={<Dashboard isLogin={isLogin} setIsLogin={setIsLogin}/>} /> */}
      </Routes>
    </div>
  );
}

export default App;
