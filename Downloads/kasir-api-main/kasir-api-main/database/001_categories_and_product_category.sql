-- 1) Categories table
CREATE TABLE IF NOT EXISTS categories (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL UNIQUE,
  description TEXT NOT NULL DEFAULT ''
);

-- 2) Add category_id to products
ALTER TABLE products
  ADD COLUMN IF NOT EXISTS category_id INT;

-- 3) Create a default category and backfill existing products (optional but recommended)
INSERT INTO categories (name, description)
VALUES ('Uncategorized', '')
ON CONFLICT (name) DO NOTHING;

UPDATE products
SET category_id = (SELECT id FROM categories WHERE name = 'Uncategorized' LIMIT 1)
WHERE category_id IS NULL;

-- 4) Enforce relationship (after backfill)
ALTER TABLE products
  ALTER COLUMN category_id SET NOT NULL;

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM pg_constraint
    WHERE conname = 'fk_products_category'
  ) THEN
    ALTER TABLE products
      ADD CONSTRAINT fk_products_category
      FOREIGN KEY (category_id)
      REFERENCES categories(id)
      ON UPDATE CASCADE
      ON DELETE RESTRICT;
  END IF;
END
$$;

CREATE INDEX IF NOT EXISTS idx_products_category_id ON products(category_id);
