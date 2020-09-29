import React from "react";

import { Slip } from "../types";

type SlipsProps = {
  slips: Slip[];
};
const Slips = ({ slips }: SlipsProps) => {
  return (
    <table className="table">
      <thead>
        <tr>
          <td>Nome</td>
          <td>Prezzo</td>
          <td></td>
        </tr>
      </thead>
      <tbody>
        {slips.map((s) => (
          <SlipRow key={s.id} slip={s} />
        ))}
      </tbody>
      <tfoot></tfoot>
    </table>
  );
};

type SlipRowProps = {
  slip: Slip;
};

const SlipRow = ({ slip }: SlipRowProps) => {
  console.log(slip);
  return (
    <tr>
      <td>{slip.name}</td>
      <td>{slip.rate}</td>
      <td>{slip.created_at}</td>
      <td></td>
    </tr>
  );
};

export default Slips;
