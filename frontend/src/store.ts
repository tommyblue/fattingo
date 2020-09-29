import { computed, decorate, observable } from "mobx";
import { Customer, Slip } from "./types";

class Store {
  _customers: Customer[] = [];
  _slips: { [customerId: string]: Slip[] } = {};

  get Customers(): Customer[] {
    return this._customers;
  }

  set Customers(c: Customer[]) {
    this._customers = c;
  }

  Slips(customerId: number): Slip[] {
    if (customerId in this._slips) {
      return this._slips[customerId];
    }
    return [];
  }

  ActiveSlips(customerId: number): Slip[] {
    if (customerId in this._slips) {
      return this._slips[customerId].filter(
        (s) => s.invoice_id === null && s.invoice_project_id === null
      );
    }
    return [];
  }

  AddSlips(customerId: number, s: Slip[]) {
    this._slips[customerId] = s;
  }

  AddCustomer(c: Customer) {
    this._customers.push(c);
  }

  GetCustomer(id: number): Customer | null {
    let res: Customer | null = null;
    this._customers.forEach((c) => {
      if (c.id === id) {
        res = c;
      }
    });
    return res;
  }
}

decorate(Store, {
  _customers: observable,
  _slips: observable,
  Customers: computed,
});

export default Store;
