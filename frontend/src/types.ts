export interface Customer {
  id: number;
  title: string;
  name: string;
  surname: string;
  address: string;
  zip_code: string;
  town: string;
  province: string;
  country: string;
  tax_code: string;
  vat: string;
  info: string;
}

export interface Slip {
  id: number;
  customer_id: number;
  invoice_id: number;
  invoice_project_id: number;
  name: string;
  rate: number;
  created_at: string;
  updated_at: string;
}

export interface CustomerInfo {
  Customer: Customer;
  Slips: Slip[];
}
