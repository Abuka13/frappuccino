CREATE TYPE order_status AS ENUM ('open', 'closed');
CREATE TYPE payment_method AS ENUM ('cash', 'card', 'kaspi_qr');
CREATE TYPE item_size AS ENUM ('small', 'medium', 'large');
CREATE TYPE transaction_type AS ENUM ('added', 'written off', 'sale', 'created');

CREATE TABLE IF NOT EXISTS inventory (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    stock NUMERIC(10, 2) NOT NULL,
    price NUMERIC(10, 2) NOT NULL CHECK (price >= 0),
    unit_type TEXT NOT NULL,
    last_updated TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS customers (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE,
    preferences JSONB
);

CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    customer_id INT NOT NULL REFERENCES customers(id),
    total_amount NUMERIC(10, 2) NOT NULL CHECK (total_amount >= 0),
    status order_status NOT NULL DEFAULT 'open',
    special_instructions JSONB,
    payment_method payment_method NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS menu_items (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    price NUMERIC(10, 2) NOT NULL CHECK (price >= 0),
    allergens TEXT[],
    category TEXT,
    size item_size NOT NULL,
    CONSTRAINT unique_menu_item_size UNIQUE (name, size)
);

CREATE TABLE IF NOT EXISTS order_items (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    menu_item_id TEXT NOT NULL REFERENCES menu_items(id) ON DELETE CASCADE,
    quantity NUMERIC NOT NULL CHECK (quantity > 0),
    price_at_order NUMERIC(10, 2) NOT NULL CHECK (price_at_order >= 0)
);

CREATE TABLE IF NOT EXISTS menu_item_ingredients (
    id SERIAL PRIMARY KEY,
    menu_item_id TEXT NOT NULL REFERENCES menu_items(id) ON DELETE CASCADE,
    ingredient_id TEXT NOT NULL REFERENCES inventory(id) ON DELETE CASCADE,
    quantity NUMERIC(10, 5) NOT NULL CHECK (quantity > 0),
    CONSTRAINT unique_menu_item_ingredient UNIQUE (menu_item_id, ingredient_id)
);

CREATE TABLE IF NOT EXISTS order_status_history (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    previous_status order_status NOT NULL,
    new_status order_status NOT NULL,
    changed_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS price_history (
    id SERIAL PRIMARY KEY,
    menu_item_id TEXT NOT NULL REFERENCES menu_items(id) ON DELETE CASCADE,
    old_price NUMERIC(10, 2) NOT NULL CHECK (old_price >= 0),
    new_price NUMERIC(10, 2) NOT NULL CHECK (new_price >= 0),
    changed_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS inventory_transactions (
    id SERIAL PRIMARY KEY,
    inventory_id TEXT NOT NULL REFERENCES inventory(id) ON DELETE CASCADE,
    change_amount NUMERIC NOT NULL,
    transaction_type transaction_type NOT NULL,
    changed_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_orders_customer_id ON orders(customer_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_order_items_menu_item_id ON order_items(menu_item_id);
CREATE INDEX idx_menu_items_name ON menu_items(name);
CREATE INDEX idx_inventory_name ON inventory(name);
CREATE INDEX idx_inventory_price ON inventory(price);
CREATE INDEX idx_inventory_stock_level ON inventory(stock);
CREATE INDEX idx_order_status_history_order_id ON order_status_history(order_id);
CREATE INDEX idx_price_history_menu_item_id ON price_history(menu_item_id);
CREATE INDEX idx_menu_items_description ON menu_items (description);
CREATE INDEX idx_customers_name ON customers (name);

-- Insert mock data into the inventory table
INSERT INTO inventory (id, name, stock, unit_type, price) VALUES
('1', 'Espresso Beans', 150, 'kg', 12.0),
('2', 'Milk', 100, 'liters', 2.0),
('3', 'Chocolate Syrup', 50, 'liters', 8.0),
('4', 'Bread', 200, 'loafs', 1.2),
('5', 'Cheese', 100, 'kg', 10.0),
('6', 'Lettuce', 50, 'kg', 3.0),
('7', 'Chicken', 80, 'kg', 12.0),
('8', 'Tomatoes', 40, 'kg', 4.0),
('9', 'Avocados', 60, 'kg', 5.0),
('10', 'Olive Oil', 20, 'liters', 10.0),
('11', 'Butter', 30, 'kg', 7.0),
('12', 'Spinach', 45, 'kg', 3.0),
('13', 'Pasta', 200, 'kg', 2.0),
('14', 'Fruit Mixture', 100, 'kg', 5.5),
('15', 'Yogurt', 50, 'liters', 4.0),
('16', 'Croissants', 120, 'pieces', 1.5),
('17', 'Bagels', 180, 'pieces', 1.8),
('18', 'Cream', 30, 'liters', 3.0),
('19', 'Sugar', 150, 'kg', 1.2),
('20', 'Coffee Cups', 1000, 'pieces', 0.05),
('21', 'Flour', 300, 'kg', 1.0),
('22', 'Ham', 70, 'kg', 8.0);

INSERT INTO menu_items (id, name, description, price, allergens, category, size) VALUES
('1', 'Cappuccino', 'Espresso with steamed milk and thick foam', 4.00, ARRAY['coffee', 'milk'], 'Beverage', 'medium'),
('2', 'Americano', 'Espresso diluted with hot water', 3.50, ARRAY['coffee'], 'Beverage', 'medium'),
('3', 'Flat White', 'Espresso with smooth steamed milk', 4.20, ARRAY['coffee', 'milk'], 'Beverage', 'medium'),
('4', 'Cheese Croissant', 'Croissant filled with cheese', 3.00, ARRAY['gluten', 'dairy'], 'Pastry', 'small'),
('5', 'Chocolate Croissant', 'Croissant filled with chocolate', 3.50, ARRAY['gluten', 'dairy'], 'Pastry', 'small'),
('6', 'Muffin', 'Soft baked muffin', 2.80, ARRAY['gluten', 'dairy'], 'Pastry', 'small'),
('7', 'Sandwich', 'Ham and cheese sandwich', 5.50, ARRAY['gluten', 'dairy', 'meat'], 'Food', 'medium'),
('8', 'Espresso', 'Strong black coffee brewed by forcing steam through finely ground coffee beans', 2.50, ARRAY['coffee'], 'Beverage', 'small'),
('9', 'Latte', 'Espresso with steamed milk and a light layer of foam', 3.80, ARRAY['coffee', 'milk'], 'Beverage', 'medium'),
('10', 'Mocha', 'Espresso with chocolate syrup, steamed milk, and whipped cream', 5.00, ARRAY['coffee', 'milk', 'chocolate'], 'Beverage', 'medium'),
('11', 'Grilled Cheese Sandwich', 'Cheese sandwich with toasted bread', 5.80, ARRAY['gluten', 'dairy'], 'Food', 'medium'),
('12', 'Chicken Salad', 'Fresh salad with grilled chicken and dressing', 7.00, ARRAY['meat', 'dairy', 'gluten'], 'Food', 'large'),
('13', 'Pasta Primavera', 'Pasta with fresh vegetables in a light sauce', 9.00, ARRAY['gluten', 'dairy'], 'Food', 'large'),
('14', 'Avocado Toast', 'Toasted bread with mashed avocado, sprinkled with chili flakes', 6.00, ARRAY['gluten', 'vegan'], 'Food', 'medium'),
('15', 'Mixed Berry Smoothie', 'Blended mixed berries with yogurt', 4.50, ARRAY['dairy', 'fruit'], 'Beverage', 'large'),
('16', 'Croissant', 'Flaky, buttery pastry', 2.00, ARRAY['gluten', 'dairy'], 'Pastry', 'small'),
('17', 'Bagel with Cream Cheese', 'Soft bagel with a layer of cream cheese', 2.50, ARRAY['gluten', 'dairy'], 'Pastry', 'small');

INSERT INTO menu_item_ingredients (menu_item_id, ingredient_id, quantity) VALUES
('8', '1', 0.02000),
('1', '1', 0.02000),
('1', '2', 0.05000),
('9', '1', 0.02000),
('9', '2', 0.08000),
('2', '1', 0.03000),
('3', '1', 0.02000),
('3', '2', 0.06000),
('4', '4', 0.10000),
('4', '11', 0.05000),
('4', '5', 0.05000),
('5', '4', 0.10000),
('5', '11', 0.05000),
('5', '3', 0.05000),
('6', '4', 0.10000),
('6', '11', 0.05000),
('6', '19', 0.05000),
('17', '4', 0.12000),
('17', '11', 0.03000),
('17', '5', 0.05000),
('7', '4', 0.15000),
('7', '5', 0.05000),
('7', '22', 0.08000);

-- Insert customers (now with 30 records to match all orders)
INSERT INTO customers (name, email, preferences) VALUES
('John Smith', 'john_smith@gmail.com', '{"note":"subscribe_to_newsletters"}'),
('Emily Johnson', 'emily_johnson@gmail.com', '{"note":"prefers_clothing_discounts"}'),
('Michael Williams', 'michael_williams@gmail.com', '{"note":"not_interested_in_ads"}'),
('Sarah Brown', 'sarah_brown@gmail.com', '{"note":"interested_in_electronics_promotions"}'),
('David Jones', 'david_jones@gmail.com', '{"note":"wants_product_updates"}'),
('Olivia Garcia', 'olivia_garcia@gmail.com', '{"note":"prefers_home_goods"}'),
('James Martinez', 'james_martinez@gmail.com', '{"note":"interested_in_eco_friendly_products"}'),
('Sophia Rodriguez', 'sophia_rodriguez@gmail.com', '{"note":"interested_in_new_books"}'),
('Daniel Wilson', 'daniel_wilson@gmail.com', '{"note":"wants_travel_promotions"}'),
('Isabella Moore', 'isabella_moore@gmail.com', '{"note":"prefers_cosmetics_discounts"}'),
('William Taylor', 'william_taylor@gmail.com', '{"note":"interested_in_sports_and_fitness"}'),
('Charlotte Anderson', 'charlotte_anderson@gmail.com', '{"note":"not_interested_in_newsletters"}'),
('Lucas Thomas', 'lucas_thomas@gmail.com', '{"note":"interested_in_pets"}'),
('Mia Jackson', 'mia_jackson@gmail.com', '{"note":"prefers_baby_products"}'),
('Henry White', 'henry_white@gmail.com', '{"note":"looking_for_travel_deals"}'),
('Emma Harris', 'emma_harris@gmail.com', '{"note":"interested_in_fashion"}'),
('Alexander Clark', 'alexander_clark@gmail.com', '{"note":"prefers_tech_products"}'),
('Ava Lewis', 'ava_lewis@gmail.com', '{"note":"wants_food_delivery_deals"}'),
('Benjamin Walker', 'benjamin_walker@gmail.com', '{"note":"interested_in_diy"}'),
('Chloe Hall', 'chloe_hall@gmail.com', '{"note":"prefers_beauty_products"}'),
('Jacob Allen', 'jacob_allen@gmail.com', '{"note":"interested_in_gaming"}'),
('Abigail Young', 'abigail_young@gmail.com', '{"note":"wants_book_recommendations"}'),
('Matthew Hernandez', 'matthew_hernandez@gmail.com', '{"note":"interested_in_fitness"}'),
('Elizabeth King', 'elizabeth_king@gmail.com', '{"note":"prefers_home_decor"}'),
('Ethan Wright', 'ethan_wright@gmail.com', '{"note":"interested_in_cars"}'),
('Sofia Lopez', 'sofia_lopez@gmail.com', '{"note":"wants_recipe_ideas"}'),
('Andrew Hill', 'andrew_hill@gmail.com', '{"note":"interested_in_photography"}'),
('Madison Scott', 'madison_scott@gmail.com', '{"note":"prefers_pet_products"}'),
('Joshua Green', 'joshua_green@gmail.com', '{"note":"interested_in_music"}'),
('Victoria Adams', 'victoria_adams@gmail.com', '{"note":"wants_travel_tips"}');

-- Now all orders can be inserted without foreign key violations
INSERT INTO orders (customer_id, total_amount, status, special_instructions, payment_method, created_at, updated_at) VALUES
(1, 7.80, 'open', '{"note":"Add extra milk"}', 'card', '2024-01-10 08:45:00', '2024-01-10 08:50:00'),
(2, 12.00, 'closed', '{"note":"Extra cheese"}', 'cash', '2024-01-12 12:30:00', '2024-01-12 12:35:00'),
(3, 8.50, 'open', '{"note":"No sugar"}', 'card', '2024-01-14 15:15:00', '2024-01-14 15:20:00'),
(4, 15.00, 'closed', '{"note":"Add extra chicken"}', 'cash', '2024-01-16 13:25:00', '2024-01-16 13:30:00'),
(5, 10.00, 'open', '{"note":"Extra avocado"}', 'cash', '2024-01-18 16:00:00', '2024-01-18 16:05:00'),
(6, 6.80, 'closed', '{"note":"Add cream"}', 'card', '2024-01-20 18:00:00', '2024-01-20 18:05:00'),
(7, 5.20, 'open', '{"note":"No tomatoes"}', 'cash', '2024-01-22 10:45:00', '2024-01-22 10:50:00'),
(8, 20.00, 'closed', '{"note":"Add extra shot of espresso"}', 'card', '2024-01-25 11:00:00', '2024-01-25 11:05:00'),
(9, 14.00, 'open', '{"note":"Spicy chicken"}', 'card', '2024-01-27 14:30:00', '2024-01-27 14:35:00'),
(10, 9.50, 'closed', '{"note":"Add extra cinnamon"}', 'cash', '2024-01-29 17:00:00', '2024-01-29 17:05:00'),
(11, 11.20, 'open', '{"note":"No onions"}', 'card', '2024-02-01 09:00:00', '2024-02-01 09:05:00'),
(12, 7.90, 'closed', '{"note":"Extra fruit"}', 'card', '2024-02-02 14:10:00', '2024-02-02 14:15:00'),
(13, 5.60, 'open', '{"note":"No dairy"}', 'cash', '2024-02-04 18:45:00', '2024-02-04 18:50:00'),
(14, 13.00, 'closed', '{"note":"Extra avocado"}', 'card', '2024-02-06 12:30:00', '2024-02-06 12:35:00'),
(15, 9.80, 'open', '{"note":"Add extra sugar"}', 'cash', '2024-02-08 10:00:00', '2024-02-08 10:05:00'),
(16, 12.50, 'closed', '{"note":"No cream"}', 'card', '2024-02-10 16:30:00', '2024-02-10 16:35:00'),
(17, 8.60, 'open', '{"note":"Add extra toast"}', 'cash', '2024-02-12 14:00:00', '2024-02-12 14:05:00'),
(18, 7.30, 'closed', '{"note":"Add extra yogurt"}', 'card', '2024-02-14 13:00:00', '2024-02-14 13:05:00'),
(19, 18.20, 'open', '{"note":"No onions, extra cheese"}', 'card', '2024-02-16 15:30:00', '2024-02-16 15:35:00'),
(20, 14.80, 'closed', '{"note":"Extra sauce"}', 'cash', '2024-02-18 10:45:00', '2024-02-18 10:50:00'),
(21, 16.00, 'open', '{"note":"More tomatoes"}', 'card', '2024-02-20 08:30:00', '2024-02-20 08:35:00'),
(22, 13.50, 'closed', '{"note":"Spicy salsa"}', 'cash', '2024-02-22 17:00:00', '2024-02-22 17:05:00'),
(23, 9.00, 'open', '{"note":"No cream"}', 'card', '2024-02-24 14:30:00', '2024-02-24 14:35:00'),
(24, 7.10, 'closed', '{"note":"Extra cheese"}', 'cash', '2024-02-26 10:00:00', '2024-02-26 10:05:00'),
(25, 15.30, 'open', '{"note":"Extra avocado"}', 'card', '2024-02-28 13:45:00', '2024-02-28 13:50:00'),
(26, 8.90, 'closed', '{"note":"No spices"}', 'cash', '2024-03-01 17:30:00', '2024-03-01 17:35:00'),
(27, 6.40, 'open', '{"note":"More lettuce"}', 'card', '2024-03-03 09:00:00', '2024-03-03 09:05:00'),
(28, 11.70, 'closed', '{"note":"No sugar"}', 'card', '2024-03-05 16:00:00', '2024-03-05 16:05:00'),
(29, 10.50, 'open', '{"note":"Less salt"}', 'cash', '2024-03-07 10:30:00', '2024-03-07 10:35:00'),
(30, 8.80, 'closed', '{"note":"Extra cinnamon"}', 'card', '2024-03-09 14:15:00', '2024-03-09 14:20:00');

INSERT INTO order_items (order_id, menu_item_id, quantity, price_at_order) VALUES
(1, '8', 2, 3.50),
(1, '4', 1, 2.50),
(2, '1', 1, 4.50),
(2, '6', 2, 2.80),
(3, '9', 1, 4.00),
(3, '5', 1, 3.00),
(4, '2', 2, 3.00),
(4, '17', 1, 2.60),
(5, '3', 1, 4.20),
(5, '7', 1, 5.50),
(6, '4', 1, 2.50),
(6, '6', 2, 2.80),
(7, '9', 1, 4.00),
(7, '7', 1, 5.50),
(8, '2', 1, 3.00),
(8, '6', 2, 2.80),
(9, '8', 1, 3.50),
(9, '5', 1, 3.00),
(10, '1', 1, 4.50),
(10, '4', 1, 2.50),
(16, '8', 2, 3.50),
(16, '4', 1, 2.50),
(17, '9', 1, 4.00),
(17, '6', 2, 2.80),
(18, '2', 1, 3.00),
(18, '17', 1, 2.60),
(19, '3', 1, 4.20),
(19, '7', 1, 5.50),
(20, '8', 1, 3.50),
(20, '4', 1, 2.50),
(21, '1', 1, 4.50),
(21, '4', 1, 3.00),
(22, '9', 1, 4.00),
(22, '6', 2, 2.80),
(23, '8', 2, 3.50),
(23, '7', 1, 5.50),
(24, '2', 1, 3.00),
(24, '6', 2, 2.80),
(25, '9', 1, 4.00),
(25, '4', 1, 2.50),
(26, '3', 1, 4.20),
(26, '7', 1, 5.50),
(27, '1', 1, 4.50),
(27, '4', 1, 2.50),
(28, '9', 1, 4.00),
(28, '7', 1, 5.50),
(29, '8', 1, 3.50),
(29, '5', 1, 3.00),
(30, '1', 1, 4.50),
(30, '4', 1, 2.50);

INSERT INTO inventory_transactions (inventory_id, change_amount, transaction_type, changed_at) VALUES
('1', -1.0, 'sale', '2024-01-10'),
('2', -2.0, 'sale', '2024-01-12'),
('3', -0.5, 'sale', '2024-01-14'),
('4', -1.0, 'sale', '2024-01-16'),
('5', -2.0, 'sale', '2024-01-18'),
('6', -0.2, 'sale', '2024-01-20'),
('7', -1.0, 'sale', '2024-01-22'),
('8', -0.5, 'sale', '2024-01-25'),
('9', -3.0, 'sale', '2024-01-27'),
('10', -1.0, 'sale', '2024-01-29'),
('11', -0.5, 'sale', '2024-02-01'),
('12', -2.5, 'sale', '2024-02-02'),
('13', -0.2, 'sale', '2024-02-04'),
('14', -1.5, 'sale', '2024-02-06'),
('15', -2.5, 'sale', '2024-02-08'),
('16', -0.1, 'sale', '2024-02-10'),
('17', -1.5, 'sale', '2024-02-12'),
('18', -0.3, 'sale', '2024-02-14'),
('19', -4.0, 'sale', '2024-02-16'),
('20', -1.2, 'sale', '2024-02-18');

INSERT INTO order_status_history (order_id, previous_status, new_status, changed_at) VALUES
(1, 'open', 'closed', '2024-01-10'),
(2, 'open', 'closed', '2024-01-12'),
(3, 'open', 'closed', '2024-01-14'),
(4, 'open', 'closed', '2024-01-16'),
(5, 'open', 'closed', '2024-01-18'),
(6, 'open', 'closed', '2024-01-20'),
(7, 'open', 'closed', '2024-01-22'),
(8, 'open', 'closed', '2024-01-25'),
(9, 'open', 'closed', '2024-01-27'),
(10, 'open', 'closed', '2024-01-29'),
(11, 'open', 'closed', '2024-02-01'),
(12, 'open', 'closed', '2024-02-02'),
(13, 'open', 'closed', '2024-02-04'),
(14, 'open', 'closed', '2024-02-06'),
(15, 'open', 'closed', '2024-02-08'),
(16, 'open', 'closed', '2024-02-10'),
(17, 'open', 'closed', '2024-02-12'),
(18, 'open', 'closed', '2024-02-14'),
(19, 'open', 'closed', '2024-02-16'),
(20, 'open', 'closed', '2024-02-18');

INSERT INTO price_history (menu_item_id, old_price, new_price, changed_at) VALUES
('8', 2.00, 2.50, '2024-01-01'),
('9', 3.50, 3.80, '2024-01-01'),
('10', 4.50, 5.00, '2024-01-01'),
('11', 5.00, 5.80, '2024-01-01'),
('12', 6.50, 7.00, '2024-01-01'),
('13', 8.50, 9.00, '2024-01-01'),
('14', 5.50, 6.00, '2024-01-01'),
('3', 4.00, 4.50, '2024-01-01'),
('4', 1.80, 2.00, '2024-01-01'),
('6', 2.20, 2.50, '2024-01-01'),
('11', 2.50, 3.00, '2024-02-01'),
('14', 3.80, 4.00, '2024-02-01'),
('12', 5.00, 5.50, '2024-02-01'),
('15', 5.80, 6.00, '2024-02-01'),
('4', 7.00, 7.50, '2024-02-01'),
('5', 9.00, 9.50, '2024-02-01'),
('5', 6.00, 6.20, '2024-02-01'),
('1', 4.50, 5.00, '2024-02-01'),
('5', 2.00, 2.20, '2024-02-01'),
('9', 2.50, 2.80, '2024-02-01');