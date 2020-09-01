import { computed, decorate, observable } from "mobx";
import { Customer } from "./types";

class Store {
  _customers: Customer[] = [];

  get Customers(): Customer[] {
    return this._customers;
  }

  set Customers(c: Customer[]) {
    this._customers = c;
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
  Customers: computed,
});

export default Store;
