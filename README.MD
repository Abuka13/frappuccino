Frappuccino ☕
Frappuccino is a Go-based project that refactors application to use a PostgreSQL database, replacing JSON-based storage with a scalable relational database. It focuses on CRUD operations, Entity-Relationship Diagram (ERD) design, and leverages PostgreSQL features like JSONB, arrays, ENUMs, and timestamps with time zones. The project includes core tables for orders, menu items, inventory, and more, with new endpoints for aggregation and reporting (e.g., total sales, popular items, full-text search). Containerized with Docker and Docker Compose for easy setup, it ensures data integrity with transactions and proper indexing.


📖 Table of Contents

Project Overview
Features
Technologies Used
Setup Instructions
Database Schema
Endpoints
Development Guidelines
Contributing
License


🌟 Project Overview
Frappuccino builds on the hot-coffee project by transitioning from JSON-based data storage to a PostgreSQL database. It emphasizes SQL, PostgreSQL, and database design principles while maintaining a clean, layered architecture. The project includes CRUD operations, advanced PostgreSQL features, and new reporting endpoints, all containerized for seamless deployment and testing.
Learning Objectives:

SQL and PostgreSQL proficiency
CRUD operations
ERD design
PostgreSQL advanced features (JSONB, arrays, ENUMs, timestamps)
Containerization with Docker


✨ Features

Database Migration: Rewrites hot-coffee JSON handlers to use PostgreSQL with a well-designed ERD.
Core Tables: Manages orders, menu items, inventory, order status history, price history, and inventory transactions.
Advanced PostgreSQL Features:
JSONB for customization and metadata
Arrays for categories and allergens
ENUMs for statuses and roles
Timestamps with time zones for tracking


New Endpoints:
Aggregation reports (total sales, popular items)
Full-text search across orders and menu items
Inventory leftovers with sorting and pagination
Bulk order processing with transaction support


Containerized Setup: Uses Docker and Docker Compose for easy deployment.
Data Integrity: Ensures consistency with transactions and proper indexing.


🛠 Technologies Used

Go: Backend programming language
PostgreSQL: Relational database
Docker & Docker Compose: Containerization for setup and testing
pgx: PostgreSQL driver for Go
gofumpt: Code formatting tool
Markdown: For documentation


🚀 Setup Instructions
Prerequisites

Docker and Docker Compose installed
Go (for development)
Git (to clone the repository)

Steps

Clone the Repository:
git clone https://github.com/your-username/frappuccino.git
cd frappuccino


Set Up Configuration:

Ensure the Dockerfile and docker-compose.yml files are in the root directory.
Create an init.sql file in the root directory with the database schema and mock data (see Database Schema).


Run the Application:
docker compose up

This starts the application at localhost:8080 and the PostgreSQL database with the following credentials:

Host: db
Port: 5432
User: latte
Password: latte
Database: frappuccino


Access the API:

Use tools like Postman or curl to interact with endpoints at http://localhost:8080.


Stop the Application:
docker compose down




🗄 Database Schema
The database schema is defined in the init.sql file and includes:
Core Tables

orders: Tracks order details, status, and timestamps.
order_items: Links orders to menu items with quantities and prices.
menu_items: Stores product details, categories, and pricing.
menu_item_ingredients: Junction table for recipes and inventory.
inventory: Manages ingredient stock levels.
order_status_history: Tracks order state changes.
price_history: Records menu item price changes.
inventory_transactions: Logs ingredient usage.

PostgreSQL Features

JSONB: For customization options, special instructions, and metadata.
Arrays: For item categories, allergens, and statuses.
ENUMs: For order statuses, payment methods, and staff roles.
Timestamps with Time Zone: For order dates and inventory updates.

Indexes

At least two indexes for frequently queried columns.
Full-text search indexes for search functionality.
Composite indexes for performance optimization.

Mock Data

10+ menu items (e.g., latte, muffin) with prices and categories.
20+ inventory items with varying quantities.
30+ orders in different statuses.
Price history and order status transitions for testing.


🌐 Endpoints
Orders

POST /orders: Create a new order.
GET /orders: Retrieve all orders.
GET /orders/{id}: Retrieve a specific order.
PUT /orders/{id}: Update an order.
DELETE /orders/{id}: Delete an order.
POST /orders/{id}/close: Close an order.

Menu Items

POST /menu: Add a new menu item.
GET /menu: Retrieve all menu items.
GET /menu/{id}: Retrieve a specific menu item.
PUT /menu/{id}: Update a menu item.
DELETE /menu/{id}: Delete a menu item.

Inventory Management 🛒:

Manage inventory with the following CRUD operations:

    POST /inventory: ➕ Add a new item to the inventory.

    GET /inventory: 📜 Retrieve a list of all inventory items.

    GET /inventory/{id}: 🔍 Fetch details of a specific inventory item by its ID.

    PUT /inventory/{id}: ✏️ Update the details of an existing inventory item by its ID.

    DELETE /inventory/{id}: ❌ Delete a specific inventory item from the system.
