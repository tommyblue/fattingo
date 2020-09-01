import React from "react";
import { useParams } from "react-router-dom";
import { observer } from "mobx-react";

import Store from "../store";

type CustomersProps = {
  store: Store;
};

const Pages = observer(({ store }: CustomersProps) => {
  let { customerId } = useParams<{ customerId: string }>();
  return (
    <div>
      <h3>Requested topic ID: {customerId}</h3>
      <p>Title: {store.GetCustomer(Number(customerId))?.title}</p>
    </div>
  );
});

export default Pages;
