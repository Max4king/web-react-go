import React from "react";
import { Routes, Route, Link } from "react-router-dom";
import LoginPage from "./LoginPage";
import "./App.css";
import RegisterPage from "./RegisterPage";
import { useState, useEffect } from "react";
import Dashboard from "./Dashboard";
import AdminDashboard from "./AdminDashboard";

function App() {
  const [isLoggedIn, setIsLoggedIn] = useState(
    localStorage.getItem("isLoggedIn") === "true"
  );
  const [isAdmin, setIsAdmin] = useState(
    localStorage.getItem("isAdmin") === "true"
  );
  const getAdminStatus = async () => {
    try {
      const response2 = await fetch("http://localhost:1323/api/isadmin", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${localStorage.getItem("token")}`,
        },
      });
      if (response2.ok) {
        const data2 = await response2.json(); // Parse JSON only if response is OK
        localStorage.setItem("isAdmin", data2);
        setIsAdmin(localStorage.getItem("isAdmin"));
        console.log(isAdmin);
      } else {
        throw new Error("Failed to fetch admin status");
      }
    } catch (err) {
      console.error("Fetch error:", err);
    }
  };
  useEffect(() => {
    const updateLoginStatus = () => {
      setIsLoggedIn(localStorage.getItem("isLoggedIn") === "true");
      setIsAdmin(localStorage.getItem("isAdmin") === "true");
    };
  }, []);

  useEffect(() => {
    if (isLoggedIn) {
      getAdminStatus();
    } else {
      setIsAdmin(false); // Ensure isAdmin is set to false when logged out
    }
  }, [isLoggedIn]); // Dependency array includes isLoggedIn
  

  function logout() {
    localStorage.removeItem("isLoggedIn");
    localStorage.removeItem("token");
    localStorage.removeItem("isAdmin");
    setIsLoggedIn(false);
    setIsAdmin(false);
  }

  return (
    <div>
      <nav className="navbar navbar-text bg-body-tertiary justify-content-center">
        <Link className="nav-link me-2" to="/dashboard">
          Dashboard
        </Link>
        {isLoggedIn ? (
          <>
            {isAdmin && (
              <Link className="nav-link me-2" to="/admin-dashboard">
                Admin Dashboard
              </Link>
            )}
            <a href="/" onClick={logout}>
              Logout
            </a>
          </>
        ) : (
          <>
            <Link className="nav-link me-2" to="/">
              Login
            </Link>
            <Link className="nav-link" to="/register">
              Register
            </Link>
          </>
        )}
      </nav>
      <div>
        {isLoggedIn ? (
          <p className="text-center">
            Logged in as {localStorage.getItem("name")}
          </p>
        ) : (
          <p className="text-center">Not logged in</p>
        )}
      </div>
      <Routes>
        <Route path="/" element={<LoginPage setIsLoggedIn={setIsLoggedIn} />} />
        <Route path="/register" element={<RegisterPage />} />
        <Route path="/dashboard" element={<Dashboard />} />
        <Route path="/admin-dashboard" element={<AdminDashboard />} />
      </Routes>
    </div>
  );
}

export default App;
