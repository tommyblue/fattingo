import React from "react";
import { useParams } from "react-router-dom";

export default function Customers() {
  let { customerId } = useParams<{ customerId: string }>();
  return <h3>Requested topic ID: {customerId}</h3>;
}
