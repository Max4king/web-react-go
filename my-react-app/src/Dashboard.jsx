import React, { useState, useEffect } from 'react';

function Dashboard() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [data, setData] = useState([]); // State to store the fetched data

  const fetchData = async () => {
    setLoading(true);
    setError('');
    try {
      const response = await fetch('http://localhost:1323/api/data', {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('token')}`,
        },
      });
      if (response.ok) {
        const fetchedData = await response.json();
        setData(fetchedData); // Store the data in state
      } else {
        throw new Error('Failed to fetch data');
      }
    } catch (err) {
      console.error("Fetch error:", err);
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, []); // Effect runs once after the first render

  if (!localStorage.getItem("isLoggedIn")) {
    return <div className="alert alert-warning">Please log in to view this page.</div>;
  }

  return (
    <div className="container mt-5">
      <h1 className="mb-3">Dashboard</h1>
      {loading && <div className="alert alert-info">Loading...</div>}
      {error && <div className="alert alert-danger">{error}</div>}
      {data.length > 0 ? (
        <div className="table-responsive">
          <table className="table table-hover">
            <thead className="thead-dark">
              <tr>
                <th>Email</th>
                <th>First Name</th>
                <th>Last Name</th>
                <th>Role</th>
                <th>User ID</th>
              </tr>
            </thead>
            <tbody>
              {data.map((item, index) => (
                <tr key={index}>
                  <td>{item.Email}</td>
                  <td>{item.FirstName}</td>
                  <td>{item.LastName}</td>
                  <td>{item.Role}</td>
                  <td>{item.UserID}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      ) : (
        !loading && <div className="alert alert-secondary">No data to display.</div>
      )}
    </div>
  );
}

export default Dashboard;
