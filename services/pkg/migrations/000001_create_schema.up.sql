-- Migration: create schema

-- Create extension for UUID support
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE users (
  user_id         SERIAL PRIMARY KEY,
  username        TEXT UNIQUE NOT NULL,
  email           TEXT UNIQUE NOT NULL,
  password_hash   TEXT NOT NULL,
  role            TEXT NOT NULL CHECK (role IN ('admin', 'bartender', 'guest')),
  created_at      TIMESTAMPTZ DEFAULT now(),
  updated_at      TIMESTAMPTZ DEFAULT now()
);

-- Ingredients table
CREATE TABLE ingredients (
  ingredient_id      SERIAL PRIMARY KEY,
  name               TEXT UNIQUE NOT NULL,
  category           TEXT,
  package_size_ml    NUMERIC(8,2) NOT NULL,
  package_cost_cents INTEGER NOT NULL,
  created_at         TIMESTAMPTZ DEFAULT now(),
  updated_at         TIMESTAMPTZ DEFAULT now()
);

-- Ingredient price history table
CREATE TABLE ingredient_price_history (
  ingredient_id      INT REFERENCES ingredients ON DELETE CASCADE,
  valid_from         DATE NOT NULL,
  package_cost_cents INTEGER NOT NULL,
  PRIMARY KEY (ingredient_id, valid_from)
);

-- Ingredient stock table
CREATE TABLE ingredient_stock (
  ingredient_id INT PRIMARY KEY REFERENCES ingredients ON DELETE CASCADE,
  qty_ml        NUMERIC(10,2) NOT NULL DEFAULT 0,
  updated_at    TIMESTAMPTZ   NOT NULL DEFAULT now()
);

-- Recipes table
CREATE TABLE recipes (
  recipe_id    SERIAL PRIMARY KEY,
  name         TEXT UNIQUE NOT NULL,
  method       TEXT,
  glass        TEXT,
  garnish      TEXT,
  instructions TEXT,
  created_by   INT REFERENCES users(user_id),
  created_at   TIMESTAMPTZ DEFAULT now(),
  updated_at   TIMESTAMPTZ DEFAULT now()
);

-- Recipe items table
CREATE TABLE recipe_items (
  recipe_id     INT REFERENCES recipes      ON DELETE CASCADE,
  ingredient_id INT REFERENCES ingredients  ON DELETE RESTRICT,
  amount_ml     NUMERIC(8,2) NOT NULL,
  PRIMARY KEY (recipe_id, ingredient_id)
);

-- Bartender skills table
CREATE TABLE bartender_skills (
  user_id   INT REFERENCES users(user_id) ON DELETE CASCADE,
  recipe_id INT REFERENCES recipes ON DELETE CASCADE,
  PRIMARY KEY (user_id, recipe_id)
);

-- Orders table
CREATE TABLE orders (
  order_id     SERIAL PRIMARY KEY,
  customer_id  INT REFERENCES users(user_id) ON DELETE SET NULL,
  bartender_id INT REFERENCES users(user_id) ON DELETE SET NULL,
  status       TEXT CHECK (status IN ('pending', 'accepted', 'completed', 'canceled')),
  created_at   TIMESTAMPTZ DEFAULT now(),
  updated_at   TIMESTAMPTZ DEFAULT now()
);

-- Order items table
CREATE TABLE order_items (
  order_id   INT REFERENCES orders(order_id) ON DELETE CASCADE,
  recipe_id  INT REFERENCES recipes(recipe_id) ON DELETE RESTRICT,
  quantity   INT NOT NULL DEFAULT 1,
  price_cents INT NOT NULL,
  status     TEXT CHECK (status IN ('pending', 'preparing', 'ready', 'delivered', 'canceled')),
  PRIMARY KEY (order_id, recipe_id)
);

-- Inventory transactions table
CREATE TABLE inventory_transactions (
  transaction_id SERIAL PRIMARY KEY,
  ingredient_id  INT REFERENCES ingredients(ingredient_id) ON DELETE RESTRICT,
  quantity_ml    NUMERIC(10,2) NOT NULL,
  transaction_type TEXT CHECK (transaction_type IN ('purchase', 'usage', 'waste', 'adjustment')),
  reference_id   INT,  -- Can be order_id, recipe_id, etc. depending on transaction_type
  created_by     INT REFERENCES users(user_id) ON DELETE SET NULL,
  created_at     TIMESTAMPTZ DEFAULT now()
);

-- Create view for recipe costs
CREATE VIEW recipe_costs AS
SELECT 
  r.recipe_id,
  r.name,
  SUM(ri.amount_ml * (i.package_cost_cents::decimal / i.package_size_ml)) AS ingredient_cost_cents,
  (SUM(ri.amount_ml * (i.package_cost_cents::decimal / i.package_size_ml)) * 5)::integer AS suggested_price_cents
FROM recipes r
JOIN recipe_items ri ON r.recipe_id = ri.recipe_id
JOIN ingredients i ON ri.ingredient_id = i.ingredient_id
GROUP BY r.recipe_id, r.name;

-- Create indexes for performance
CREATE INDEX idx_ingredient_category ON ingredients(category);
CREATE INDEX idx_recipe_method ON recipes(method);
CREATE INDEX idx_recipe_created_by ON recipes(created_by);
CREATE INDEX idx_orders_customer_id ON orders(customer_id);
CREATE INDEX idx_orders_bartender_id ON orders(bartender_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_inventory_transactions_ingredient_id ON inventory_transactions(ingredient_id);
CREATE INDEX idx_order_items_recipe_id ON order_items(recipe_id);

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = now();
   RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for updated_at columns
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_ingredients_updated_at BEFORE UPDATE ON ingredients
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_ingredient_stock_updated_at BEFORE UPDATE ON ingredient_stock
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_recipes_updated_at BEFORE UPDATE ON recipes
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_orders_updated_at BEFORE UPDATE ON orders
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column(); 