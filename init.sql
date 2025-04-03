CREATE TYPE order_status AS ENUM ('open', 'close');
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

-- Insert mock data into the inventory table
INSERT INTO inventory (id, name, stock, unit_type, price) VALUES
('espresso_beans', 'Espresso Beans', 150, 'kg', 12.0),
('milk', 'Milk', 100, 'liters', 2.0),
('chocolate_syrup', 'Chocolate Syrup', 50, 'liters', 8.0),
('bread', 'Bread', 200, 'loafs', 1.2),
('cheese', 'Cheese', 100, 'kg', 10.0),
('lettuce', 'Lettuce', 50, 'kg', 3.0),
('chicken', 'Chicken', 80, 'kg', 12.0),
('tomatoes', 'Tomatoes', 40, 'kg', 4.0),
('avocados', 'Avocados', 60, 'kg', 5.0),
('olive_oil', 'Olive Oil', 20, 'liters', 10.0),
('butter', 'Butter', 30, 'kg', 7.0),
('spinach', 'Spinach', 45, 'kg', 3.0),
('pasta', 'Pasta', 200, 'kg', 2.0),
('fruit_mixture', 'Fruit Mixture', 100, 'kg', 5.5),
('yogurt', 'Yogurt', 50, 'liters', 4.0),
('croissants', 'Croissants', 120, 'pieces', 1.5),
('bagels', 'Bagels', 180, 'pieces', 1.8),
('cream', 'Cream', 30, 'liters', 3.0),
('sugar', 'Sugar', 150, 'kg', 1.2),
('coffee_cups', 'Coffee Cups', 1000, 'pieces', 0.05),
('flour', 'Flour', 300, 'kg', 1.0),
('ham', 'Ham', 70, 'kg', 8.0);

INSERT INTO menu_items (id, name, description, price, allergens, category, size) VALUES
('cappuccino', 'Cappuccino', 'Espresso with steamed milk and thick foam', 4.00, ARRAY['coffee', 'milk'], 'Beverage', 'medium'),
('americano', 'Americano', 'Espresso diluted with hot water', 3.50, ARRAY['coffee'], 'Beverage', 'medium'),
('flat_white', 'Flat White', 'Espresso with smooth steamed milk', 4.20, ARRAY['coffee', 'milk'], 'Beverage', 'medium'),
('cheese_croissant', 'Cheese Croissant', 'Croissant filled with cheese', 3.00, ARRAY['gluten', 'dairy'], 'Pastry', 'small'),
('chocolate_croissant', 'Chocolate Croissant', 'Croissant filled with chocolate', 3.50, ARRAY['gluten', 'dairy'], 'Pastry', 'small'),
('muffin', 'Muffin', 'Soft baked muffin', 2.80, ARRAY['gluten', 'dairy'], 'Pastry', 'small'),
('sandwich', 'Sandwich', 'Ham and cheese sandwich', 5.50, ARRAY['gluten', 'dairy', 'meat'], 'Food', 'medium'),
('espresso', 'Espresso', 'Strong black coffee brewed by forcing steam through finely ground coffee beans', 2.50, ARRAY['coffee'], 'Beverage', 'small'),
('latte', 'Latte', 'Espresso with steamed milk and a light layer of foam', 3.80, ARRAY['coffee', 'milk'], 'Beverage', 'medium'),
('mocha', 'Mocha', 'Espresso with chocolate syrup, steamed milk, and whipped cream', 5.00, ARRAY['coffee', 'milk', 'chocolate'], 'Beverage', 'medium'),
('grilled_cheese', 'Grilled Cheese Sandwich', 'Cheese sandwich with toasted bread', 5.80, ARRAY['gluten', 'dairy'], 'Food', 'medium'),
('chicken_salad', 'Chicken Salad', 'Fresh salad with grilled chicken and dressing', 7.00, ARRAY['meat', 'dairy', 'gluten'], 'Food', 'large'),
('pasta', 'Pasta Primavera', 'Pasta with fresh vegetables in a light sauce', 9.00, ARRAY['gluten', 'dairy'], 'Food', 'large'),
('avocado_toast', 'Avocado Toast', 'Toasted bread with mashed avocado, sprinkled with chili flakes', 6.00, ARRAY['gluten', 'vegan'], 'Food', 'medium'),
('smoothie', 'Mixed Berry Smoothie', 'Blended mixed berries with yogurt', 4.50, ARRAY['dairy', 'fruit'], 'Beverage', 'large'),
('croissant', 'Croissant', 'Flaky, buttery pastry', 2.00, ARRAY['gluten', 'dairy'], 'Pastry', 'small'),
('bagel', 'Bagel with Cream Cheese', 'Soft bagel with a layer of cream cheese', 2.50, ARRAY['gluten', 'dairy'], 'Pastry', 'small');

INSERT INTO menu_item_ingredients (menu_item_id, ingredient_id, quantity) VALUES
('espresso', 'espresso_beans', 0.02000),
('cappuccino', 'espresso_beans', 0.02000),
('cappuccino', 'milk', 0.05000),
('latte', 'espresso_beans', 0.02000),
('latte', 'milk', 0.08000),
('americano', 'espresso_beans', 0.03000),
('flat_white', 'espresso_beans', 0.02000),
('flat_white', 'milk', 0.06000),
('cheese_croissant', 'bread', 0.10000),
('cheese_croissant', 'butter', 0.05000),
('cheese_croissant', 'cheese', 0.05000),
('chocolate_croissant', 'bread', 0.10000),
('chocolate_croissant', 'butter', 0.05000),
('chocolate_croissant', 'chocolate_syrup', 0.05000),
('muffin', 'bread', 0.10000),
('muffin', 'butter', 0.05000),
('muffin', 'sugar', 0.05000),
('bagel', 'bread', 0.12000),
('bagel', 'butter', 0.03000),
('bagel', 'cheese', 0.05000),
('sandwich', 'bread', 0.15000),
('sandwich', 'cheese', 0.05000),
('sandwich', 'ham', 0.08000);

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
('Henry White', 'henry_white@gmail.com', '{"note":"looking_for_travel_deals"}');

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
(1, 'espresso', 2, 3.50),
(1, 'cheese_croissant', 1, 2.50),
(2, 'cappuccino', 1, 4.50),
(2, 'muffin', 2, 2.80),
(3, 'latte', 1, 4.00),
(3, 'chocolate_croissant', 1, 3.00),
(4, 'americano', 2, 3.00),
(4, 'bagel', 1, 2.60),
(5, 'flat_white', 1, 4.20),
(5, 'sandwich', 1, 5.50),
(6, 'cheese_croissant', 1, 2.50),
(6, 'muffin', 2, 2.80),
(7, 'latte', 1, 4.00),
(7, 'sandwich', 1, 5.50),
(8, 'americano', 1, 3.00),
(8, 'muffin', 2, 2.80),
(9, 'espresso', 1, 3.50),
(9, 'chocolate_croissant', 1, 3.00),
(10, 'cappuccino', 1, 4.50),
(10, 'cheese_croissant', 1, 2.50),
(16, 'espresso', 2, 3.50),
(16, 'cheese_croissant', 1, 2.50),
(17, 'latte', 1, 4.00),
(17, 'muffin', 2, 2.80),
(18, 'americano', 1, 3.00),
(18, 'bagel', 1, 2.60),
(19, 'flat_white', 1, 4.20),
(19, 'sandwich', 1, 5.50),
(20, 'espresso', 1, 3.50),
(20, 'cheese_croissant', 1, 2.50),
(21, 'cappuccino', 1, 4.50),
(21, 'chocolate_croissant', 1, 3.00),
(22, 'latte', 1, 4.00),
(22, 'muffin', 2, 2.80),
(23, 'espresso', 2, 3.50),
(23, 'sandwich', 1, 5.50),
(24, 'americano', 1, 3.00),
(24, 'muffin', 2, 2.80),
(25, 'latte', 1, 4.00),
(25, 'cheese_croissant', 1, 2.50),
(26, 'flat_white', 1, 4.20),
(26, 'sandwich', 1, 5.50),
(27, 'cappuccino', 1, 4.50),
(27, 'cheese_croissant', 1, 2.50),
(28, 'latte', 1, 4.00),
(28, 'sandwich', 1, 5.50),
(29, 'espresso', 1, 3.50),
(29, 'chocolate_croissant', 1, 3.00),
(30, 'cappuccino', 1, 4.50),
(30, 'cheese_croissant', 1, 2.50);

INSERT INTO inventory_transactions (inventory_id, change_amount, transaction_type, changed_at) VALUES
('espresso_beans', -1.0, 'sale', '2024-01-10'),
('bread', -2.0, 'sale', '2024-01-12'),
('chocolate_syrup', -0.5, 'sale', '2024-01-14'),
('cheese', -1.0, 'sale', '2024-01-16'),
('chicken', -2.0, 'sale', '2024-01-18'),
('avocados', -0.2, 'sale', '2024-01-20'),
('pasta', -1.0, 'sale', '2024-01-22'),
('cream', -0.5, 'sale', '2024-01-25'),
('croissants', -3.0, 'sale', '2024-01-27'),
('milk', -1.0, 'sale', '2024-01-29'),
('espresso_beans', -0.5, 'sale', '2024-02-01'),
('bread', -2.5, 'sale', '2024-02-02'),
('chocolate_syrup', -0.2, 'sale', '2024-02-04'),
('cheese', -1.5, 'sale', '2024-02-06'),
('chicken', -2.5, 'sale', '2024-02-08'),
('avocados', -0.1, 'sale', '2024-02-10'),
('pasta', -1.5, 'sale', '2024-02-12'),
('cream', -0.3, 'sale', '2024-02-14'),
('croissants', -4.0, 'sale', '2024-02-16'),
('milk', -1.2, 'sale', '2024-02-18');

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
('espresso', 2.00, 2.50, '2024-01-01'),
('latte', 3.50, 3.80, '2024-01-01'),
('mocha', 4.50, 5.00, '2024-01-01'),
('grilled_cheese', 5.00, 5.80, '2024-01-01'),
('chicken_salad', 6.50, 7.00, '2024-01-01'),
('pasta', 8.50, 9.00, '2024-01-01'),
('avocado_toast', 5.50, 6.00, '2024-01-01'),
('smoothie', 4.00, 4.50, '2024-01-01'),
('croissant', 1.80, 2.00, '2024-01-01'),
('bagel', 2.20, 2.50, '2024-01-01'),
('espresso', 2.50, 3.00, '2024-02-01'),
('latte', 3.80, 4.00, '2024-02-01'),
('mocha', 5.00, 5.50, '2024-02-01'),
('grilled_cheese', 5.80, 6.00, '2024-02-01'),
('chicken_salad', 7.00, 7.50, '2024-02-01'),
('pasta', 9.00, 9.50, '2024-02-01'),
('avocado_toast', 6.00, 6.20, '2024-02-01'),
('smoothie', 4.50, 5.00, '2024-02-01'),
('croissant', 2.00, 2.20, '2024-02-01'),
('bagel', 2.50, 2.80, '2024-02-01');