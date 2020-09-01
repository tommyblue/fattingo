import React from "react";
import { Switch, Route } from "react-router-dom";
import { observer } from "mobx-react";

import Customers from "./Pages/Customers";
import Store from "./store";

type PagesProps = {
  store: Store;
};

const Pages = observer(({ store }: PagesProps) => {
  return (
    <Switch>
      <Route path="/customers/:customerId">
        <Customers store={store} />
      </Route>
    </Switch>
  );
});

export default Pages;
