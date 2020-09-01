import React from "react";
import { Switch, Route } from "react-router-dom";

import Customers from "./Pages/Customers";

export default function () {
  return (
    <Switch>
      <Route path="/customers/:customerId">
        <Customers />
      </Route>
    </Switch>
  );
}
