import React from 'react';
import ReactDOM from 'react-dom';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import LoginPage from './LoginPage';
import './App.css'
import RegisterPage from './RegisterPage';
// import Dashboard from './Dashboard';

function App() {
  return (
    <div>
    <nav className='navbar navbar-text bg-body-tertiary justify-content-center'>
        <Link className='nav-link me-2' to="/dashboard">Dashboard</Link>
        <Link className='nav-link me-2' to="/">Login</Link>
        <Link className='nav-link' to="/register">Register</Link>
    </nav>
    <Routes>
        <Route path="/" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />
        {/* <Route path="/dashboard" element={<Dashboard />} /> */}
      </Routes>
    </div>
  );
}

export default App;
