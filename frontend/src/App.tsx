import React, { useState, useEffect } from "react";
import "./App.scss";
import { Customer } from "./types";
import Page from "./Page";
import Sidebar from "./Sidebar";

function App() {
  const [customers, setCustomers] = useState<Customer[]>([]);

  useEffect(() => {
    async function fetchData() {
      const response = await fetch("http://localhost:5000/api/v1/customers");
      const result = await response.json();
      setCustomers(result);
    }
    fetchData();
  }, []);

  return (
    <div className="columns">
      <div className="column is-3">
        <Sidebar customers={customers} />
      </div>
      <div className="column">
        <Page />
      </div>
    </div>
  );
}

export default App;
