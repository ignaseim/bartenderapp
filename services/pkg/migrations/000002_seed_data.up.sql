-- Seed data for testing and development

-- Insert sample users (password is 'password' hashed)
INSERT INTO users (username, email, password_hash, role) VALUES 
('admin', 'admin@bartenderapp.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'admin'),
('bartender1', 'bartender1@bartenderapp.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'bartender'),
('bartender2', 'bartender2@bartenderapp.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'bartender');

-- Insert sample ingredients
INSERT INTO ingredients (name, category, package_size_ml, package_cost_cents) VALUES
('Vodka', 'Spirit', 750, 1500),
('Gin', 'Spirit', 750, 2000),
('Rum', 'Spirit', 750, 1800),
('Tequila', 'Spirit', 750, 2500),
('Whiskey', 'Spirit', 750, 3000),
('Triple Sec', 'Liqueur', 750, 1500),
('Simple Syrup', 'Syrup', 500, 500),
('Lime Juice', 'Juice', 500, 600),
('Lemon Juice', 'Juice', 500, 600),
('Orange Juice', 'Juice', 1000, 800),
('Cranberry Juice', 'Juice', 1000, 700),
('Soda Water', 'Mixer', 1000, 300),
('Tonic Water', 'Mixer', 1000, 400),
('Coca Cola', 'Mixer', 1000, 500),
('Mint', 'Garnish', 100, 300),
('Olives', 'Garnish', 200, 400),
('Angostura Bitters', 'Bitters', 200, 1200),
('Vermouth', 'Fortified Wine', 750, 1500),
('Coffee Liqueur', 'Liqueur', 750, 2000),
('Cream', 'Dairy', 500, 600);

-- Set initial stock for ingredients
INSERT INTO ingredient_stock (ingredient_id, qty_ml)
SELECT ingredient_id, 2000 FROM ingredients;

-- Insert sample recipes
INSERT INTO recipes (name, method, glass, garnish, instructions, created_by) VALUES
('Mojito', 'Muddled', 'Highball', 'Mint Sprig', 'Muddle mint leaves with sugar and lime juice. Add rum, fill with ice and top with soda water.', 1),
('Martini', 'Stirred', 'Martini', 'Olive or Lemon Twist', 'Stir gin and vermouth with ice, strain into chilled glass.', 1),
('Margarita', 'Shaken', 'Coupe or Rocks', 'Salt Rim, Lime Wheel', 'Shake tequila, triple sec, and lime juice with ice. Strain into salt-rimmed glass.', 1),
('Old Fashioned', 'Built', 'Rocks', 'Orange Peel', 'Muddle sugar with bitters and water. Add whiskey and ice, stir.', 1),
('Cosmopolitan', 'Shaken', 'Coupe', 'Lime Wedge', 'Shake vodka, triple sec, cranberry juice, and lime juice with ice. Strain into glass.', 1),
('White Russian', 'Built', 'Rocks', 'None', 'Add vodka and coffee liqueur to glass with ice. Top with cream.', 1),
('Tom Collins', 'Built', 'Highball', 'Lemon Wheel', 'Combine gin, lemon juice, and simple syrup in a glass with ice. Top with soda water.', 1),
('Daiquiri', 'Shaken', 'Coupe', 'Lime Wheel', 'Shake rum, lime juice, and simple syrup with ice. Strain into glass.', 1),
('Moscow Mule', 'Built', 'Copper Mug', 'Lime Wheel', 'Add vodka and lime juice to mug with ice. Top with ginger beer.', 1),
('Whiskey Sour', 'Shaken', 'Rocks', 'Orange Slice and Cherry', 'Shake whiskey, lemon juice, and simple syrup with ice. Strain into glass.', 1);

-- Insert recipe ingredients

-- Mojito
INSERT INTO recipe_items (recipe_id, ingredient_id, amount_ml) VALUES
(1, 3, 60),  -- Rum
(1, 8, 30),  -- Lime Juice
(1, 7, 15),  -- Simple Syrup
(1, 15, 10), -- Mint
(1, 12, 90); -- Soda Water

-- Martini
INSERT INTO recipe_items (recipe_id, ingredient_id, amount_ml) VALUES
(2, 2, 60),   -- Gin
(2, 18, 15),  -- Vermouth
(2, 16, 5);   -- Olives

-- Margarita
INSERT INTO recipe_items (recipe_id, ingredient_id, amount_ml) VALUES
(3, 4, 50),   -- Tequila
(3, 6, 30),   -- Triple Sec
(3, 8, 30);   -- Lime Juice

-- Old Fashioned
INSERT INTO recipe_items (recipe_id, ingredient_id, amount_ml) VALUES
(4, 5, 60),   -- Whiskey
(4, 7, 5),    -- Simple Syrup
(4, 17, 5);   -- Angostura Bitters

-- Cosmopolitan
INSERT INTO recipe_items (recipe_id, ingredient_id, amount_ml) VALUES
(5, 1, 45),   -- Vodka
(5, 6, 15),   -- Triple Sec
(5, 11, 30),  -- Cranberry Juice
(5, 8, 15);   -- Lime Juice

-- White Russian
INSERT INTO recipe_items (recipe_id, ingredient_id, amount_ml) VALUES
(6, 1, 50),   -- Vodka
(6, 19, 30),  -- Coffee Liqueur
(6, 20, 30);  -- Cream

-- Tom Collins
INSERT INTO recipe_items (recipe_id, ingredient_id, amount_ml) VALUES
(7, 2, 60),   -- Gin
(7, 9, 30),   -- Lemon Juice
(7, 7, 20),   -- Simple Syrup
(7, 12, 90);  -- Soda Water

-- Daiquiri
INSERT INTO recipe_items (recipe_id, ingredient_id, amount_ml) VALUES
(8, 3, 60),   -- Rum
(8, 8, 30),   -- Lime Juice
(8, 7, 15);   -- Simple Syrup

-- Moscow Mule (assume we have ginger beer as id 21)
INSERT INTO recipe_items (recipe_id, ingredient_id, amount_ml) VALUES
(9, 1, 60),   -- Vodka
(9, 8, 15);   -- Lime Juice

-- Whiskey Sour
INSERT INTO recipe_items (recipe_id, ingredient_id, amount_ml) VALUES
(10, 5, 60),  -- Whiskey
(10, 9, 30),  -- Lemon Juice
(10, 7, 20);  -- Simple Syrup

-- Define bartender skills
INSERT INTO bartender_skills (user_id, recipe_id)
VALUES 
(2, 1), (2, 2), (2, 3), (2, 5), (2, 8), -- Bartender 1 skills
(3, 1), (3, 4), (3, 6), (3, 7), (3, 9), (3, 10); -- Bartender 2 skills

-- Create sample orders
INSERT INTO orders (customer_id, bartender_id, status, created_at)
VALUES 
(NULL, 2, 'completed', now() - interval '3 days'),
(NULL, 3, 'completed', now() - interval '2 days'),
(NULL, 2, 'completed', now() - interval '1 day'),
(NULL, 3, 'pending', now());

-- Create sample order items
INSERT INTO order_items (order_id, recipe_id, quantity, price_cents, status)
VALUES
(1, 1, 2, 1200, 'delivered'),
(1, 5, 1, 1400, 'delivered'),
(2, 4, 1, 1500, 'delivered'),
(2, 6, 2, 1300, 'delivered'),
(3, 3, 3, 1200, 'delivered'),
(4, 1, 1, 1200, 'pending'),
(4, 10, 2, 1400, 'pending');

-- Create sample inventory transactions
INSERT INTO inventory_transactions (ingredient_id, quantity_ml, transaction_type, reference_id, created_by)
VALUES
(3, -120, 'usage', 1, 2),   -- Rum used for order 1 (Mojito)
(8, -60, 'usage', 1, 2),    -- Lime Juice used for order 1 (Mojito)
(7, -30, 'usage', 1, 2),    -- Simple Syrup used for order 1 (Mojito)
(15, -20, 'usage', 1, 2),   -- Mint used for order 1 (Mojito)
(12, -180, 'usage', 1, 2),  -- Soda Water used for order 1 (Mojito)
(1, -45, 'usage', 1, 2),    -- Vodka used for order 1 (Cosmopolitan)
(6, -15, 'usage', 1, 2),    -- Triple Sec used for order 1 (Cosmopolitan)
(11, -30, 'usage', 1, 2),   -- Cranberry Juice used for order 1 (Cosmopolitan)
(8, -15, 'usage', 1, 2),    -- Lime Juice used for order 1 (Cosmopolitan)

(5, -60, 'usage', 2, 3),    -- Whiskey used for order 2 (Old Fashioned)
(7, -5, 'usage', 2, 3),     -- Simple Syrup used for order 2 (Old Fashioned)
(17, -5, 'usage', 2, 3),    -- Angostura Bitters used for order 2 (Old Fashioned)
(1, -100, 'usage', 2, 3),   -- Vodka used for order 2 (White Russian)
(19, -60, 'usage', 2, 3),   -- Coffee Liqueur used for order 2 (White Russian)
(20, -60, 'usage', 2, 3),   -- Cream used for order 2 (White Russian)

(4, -150, 'usage', 3, 2),   -- Tequila used for order 3 (Margarita)
(6, -90, 'usage', 3, 2),    -- Triple Sec used for order 3 (Margarita)
(8, -90, 'usage', 3, 2),    -- Lime Juice used for order 3 (Margarita)

(3, 750, 'purchase', NULL, 1),     -- Purchase of Rum
(1, 750, 'purchase', NULL, 1),     -- Purchase of Vodka
(8, 500, 'purchase', NULL, 1),     -- Purchase of Lime Juice
(11, 1000, 'purchase', NULL, 1);   -- Purchase of Cranberry Juice 