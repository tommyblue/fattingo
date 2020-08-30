CREATE TABLE customers (
    id INTEGER PRIMARY KEY,
    user_id INTEGER,
    title TEXT,
    name TEXT,
    surname TEXT,
    address TEXT,
    zip_code TEXT,
    town TEXT,
    province TEXT,
    country TEXT,
    tax_code TEXT,
    vat TEXT,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    info TEXT
);
