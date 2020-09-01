import React from "react";
import { Link } from "react-router-dom";
import { observer } from "mobx-react";

import Store from "./store";

type SidebarProps = {
  store: Store;
};

const Sidebar = observer(({ store }: SidebarProps) => {
  return (
    <aside className="menu py-4 px-4">
      <p className="menu-label">Customers</p>
      <ul className="menu-list">
        {store.Customers.slice()
          .sort((c1, c2) => (c1.title < c2.title ? -1 : 1))
          .map((customer) => (
            <li key={customer.id}>
              <Link to={`/customers/${customer.id}`}>{customer.title}</Link>
            </li>
          ))}
      </ul>
    </aside>
  );
});

export default Sidebar;
