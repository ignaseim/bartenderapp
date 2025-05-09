-- Migration: drop schema

-- Drop triggers
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_ingredients_updated_at ON ingredients;
DROP TRIGGER IF EXISTS update_ingredient_stock_updated_at ON ingredient_stock;
DROP TRIGGER IF EXISTS update_recipes_updated_at ON recipes;
DROP TRIGGER IF EXISTS update_orders_updated_at ON orders;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop view
DROP VIEW IF EXISTS recipe_costs;

-- Drop tables in reverse order to avoid foreign key constraints
DROP TABLE IF EXISTS inventory_transactions;
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS bartender_skills;
DROP TABLE IF EXISTS recipe_items;
DROP TABLE IF EXISTS recipes;
DROP TABLE IF EXISTS ingredient_stock;
DROP TABLE IF EXISTS ingredient_price_history;
DROP TABLE IF EXISTS ingredients;
DROP TABLE IF EXISTS users;

-- Drop extension
DROP EXTENSION IF EXISTS "uuid-ossp"; 