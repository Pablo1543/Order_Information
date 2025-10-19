CREATE TABLE IF NOT EXISTS orders (
  order_uid TEXT PRIMARY KEY,
  track_number TEXT,
  customer_id TEXT,
  date_created TIMESTAMP WITH TIME ZONE,
  raw JSONB NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);
