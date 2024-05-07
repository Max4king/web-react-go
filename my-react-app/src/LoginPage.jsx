import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import "./LoginPage.css";
function LoginPage({ setIsLoggedIn }) {
  const [name, setName] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [successful, setSuccessful] = useState(false);
  let navigate = useNavigate();
  const handleLogin = async (event) => {
    event.preventDefault();
    setLoading(true);
    setSuccessful(false);
    setError("");
    try {
      const response = await fetch("http://localhost:1323/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ name: name, password: password }),
      });

      if (response.ok) {
        const data = await response.json(); // Parse JSON only if response is OK
        setSuccessful(true);
        localStorage.setItem("isLoggedIn", "true");
        setIsLoggedIn(true);
        localStorage.setItem("token", data);
        localStorage.setItem('name', name);
        navigate("/dashboard");
      } else {
        const errorData = await response.text(); // Use .text() if the response might not be JSON
        throw new Error(errorData || "Failed to login");
      }
    } catch (err) {
      setError(err.message);
      console.error("Login error:", err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="login-container">
      <form className="login-form" onSubmit={handleLogin}>
        <h1>Login</h1>
        {error && <p style={{ color: "red" }}>{error}</p>}
        {successful && <p style={{ color: "green" }}>Login Successful</p>}
        <div>
          <label htmlFor="name">Name:</label>
          <input
            type="text"
            id="name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            required
          />
        </div>
        <div>
          <label htmlFor="password">Password:</label>
          <input
            type="password"
            id="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </div>
        <button type="submit" className="btn btn-primary" disabled={loading}>
          {loading ? "Logging in..." : "Login"}
        </button>
      </form>
    </div>
  );
}

export default LoginPage;
