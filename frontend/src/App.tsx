import React, { useEffect } from "react";
import { BrowserRouter as Router } from "react-router-dom";

import "./App.scss";
import Pages from "./Pages";
import Sidebar from "./Sidebar";
import Store from "./store";

const App = () => {
  const store = new Store();

  useEffect(() => {
    async function fetchData() {
      const response = await fetch("http://localhost:5000/api/v1/customers");
      const result = await response.json();
      store.Customers = result;
    }
    fetchData();
  }, [store.Customers]);

  return (
    <Router>
      <div className="columns">
        <div className="column is-3">
          <Sidebar store={store} />
        </div>
        <div className="column">
          <Pages store={store} />
        </div>
      </div>
    </Router>
  );
};

export default App;
