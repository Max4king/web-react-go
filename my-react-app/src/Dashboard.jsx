import React, { useState, useEffect } from "react";

function Dashboard() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [data, setData] = useState([]);
  const [newEmail, setNewEmail] = useState("");
  const [newFirstName, setNewFirstName] = useState("");
  const [newLastName, setNewLastName] = useState("");
  const [newRole, setNewRole] = useState("");
  const [isAdmin, setIsAdmin] = useState(false);

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
        setIsAdmin(data2);
      } else {
        throw new Error("Failed to fetch admin status");
      }
    } catch (err) {
      console.error("Fetch error:", err);
    }
  };
  
  const fetchData = async () => {
    setLoading(true);
    setError("");
    try {
      const response = await fetch("http://localhost:1323/api/data", {
        headers: {
          Authorization: `Bearer ${localStorage.getItem("token")}`,
        },
      });
      if (response.ok) {
        const fetchedData = await response.json();
        setData(fetchedData); // Store the data in state
      } else {
        throw new Error("Failed to fetch data");
      }
    } catch (err) {
      console.error("Fetch error:", err);
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    getAdminStatus();
    fetchData();
  }, []);

  const handleAddData = async (e) => {
    e.preventDefault();
    try {
      const response = await fetch("http://localhost:1323/api/new/data", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${localStorage.getItem("token")}`,
        },
        body: JSON.stringify({
          Email: newEmail,
          FirstName: newFirstName,
          LastName: newLastName,
          Role: newRole,
        }),
      });
      if (response.ok) {
        fetchData(); // Refresh data after adding
      } else {
        throw new Error("Failed to add data");
      }
    } catch (err) {
      console.error("Add data error:", err);
      setError(err.message);
    }
  };
  if (!localStorage.getItem("isLoggedIn")) {
    return (
      <div className="alert alert-warning">
        Please log in to view this page.
      </div>
    );
  }

  return (
    <div className="container mt-5">
      <h1 className="mb-3">{isAdmin && <span>Admin</span> } Dashboard</h1>
      {loading && <div className="alert alert-info">Loading...</div>}
      {error && <div className="alert alert-danger">{error}</div>}
      <form onSubmit={handleAddData}>
        <input
          type="email"
          value={newEmail}
          onChange={(e) => setNewEmail(e.target.value)}
          placeholder="Email"
          required
        />
        <input
          type="text"
          value={newFirstName}
          onChange={(e) => setNewFirstName(e.target.value)}
          placeholder="First Name"
          required
        />
        <input
          type="text"
          value={newLastName}
          onChange={(e) => setNewLastName(e.target.value)}
          placeholder="Last Name"
          required
        />
        <input
          type="text"
          value={newRole}
          onChange={(e) => setNewRole(e.target.value)}
          placeholder="Role"
          required
        />
        <button type="submit">Add Data</button>
      </form>
      {data.length > 0 ? (
        <div className="table-responsive">
          <table className="table table-hover">
            <thead className="thead-dark">
              <tr>
                <th>Email</th>
                <th>First Name</th>
                <th>Last Name</th>
                <th>Role</th>
              </tr>
            </thead>
            <tbody>
              {data.map((item, index) => (
                <tr key={index}>
                  <td>{item.Email}</td>
                  <td>{item.FirstName}</td>
                  <td>{item.LastName}</td>
                  <td>{item.Role}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      ) : (
        !loading && (
          <div className="alert alert-secondary">No data to display.</div>
        )
      )}
    </div>
  );
}

export default Dashboard;
