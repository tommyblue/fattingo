import React, { useEffect } from "react";
import { useParams } from "react-router-dom";
import { observer } from "mobx-react";
import { toJS } from "mobx";

import { CustomerInfo } from "../types";
import Store from "../store";
import Slips from "../Components/Slips";

type CustomersProps = {
  store: Store;
};

const Pages = observer(({ store }: CustomersProps) => {
  let { customerId } = useParams<{ customerId: string }>();

  useEffect(() => {
    async function fetchData() {
      const response = await fetch(
        `http://localhost:5000/api/v1/customers/${customerId}/info`
      );
      const result: CustomerInfo = await response.json();
      store.AddSlips(result.Customer.id, result.Slips);
    }
    fetchData();
  }, [customerId, store]);

  const slips = toJS(store.ActiveSlips(parseInt(customerId)));
  return (
    <div>
      <h3>Requested topic ID: {customerId}</h3>
      <p>Title: {store.GetCustomer(Number(customerId))?.title}</p>
      <Slips slips={slips} />
      <section className="section">
        <div className="container">
          <h1 className="title">Section</h1>
          <h2 className="subtitle">
            A simple container to divide your page into{" "}
            <strong>sections</strong>, like the one you're currently reading
          </h2>
        </div>
      </section>
    </div>
  );
});

export default Pages;
