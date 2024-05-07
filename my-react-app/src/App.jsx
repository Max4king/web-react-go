import React from 'react';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import LoginPage from './LoginPage';
import './App.css'
import RegisterPage from './RegisterPage';
import { useState, useEffect } from 'react';
import Dashboard from './Dashboard';


function App() {
  const [isLoggedIn, setIsLoggedIn] = useState(localStorage.getItem('isLoggedIn') === 'true');
  useEffect(() => {
    const updateLoginStatus = () => setIsLoggedIn(localStorage.getItem('isLoggedIn') === 'true');
    window.addEventListener('storage', updateLoginStatus);

    return () => window.removeEventListener('storage', updateLoginStatus);
  }, []);

  function logout() {
    localStorage.removeItem('isLoggedIn');
    localStorage.removeItem('token');
    setIsLoggedIn(false);
  }
  if (isLoggedIn) {
    return (
    <div>
    <nav className='navbar navbar-text bg-body-tertiary justify-content-center'>
        <Link className='nav-link me-2' to="/dashboard">Dashboard</Link>
        <a href='#' onClick={() => logout() }>Logout</a>
    </nav>
    <Routes>
        <Route path="/" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />
        <Route path="/dashboard" element={<Dashboard/>} />
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
        <Route path="/" element={<LoginPage/>} />
        <Route path="/register" element={<RegisterPage/>} />
        <Route path="/dashboard" element={<Dashboard/>} />
      </Routes>
    </div>
  );
}

export default App;
