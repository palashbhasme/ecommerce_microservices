CREATE DATABASE inventory;
CREATE DATABASE orders;

CREATE USER inventory_user WITH ENCRYPTED PASSWORD 'inventory_password';
CREATE USER order_user WITH ENCRYPTED PASSWORD 'order_password';

GRANT ALL PRIVILEGES ON DATABASE inventory TO inventory_user;
GRANT ALL PRIVILEGES ON DATABASE orders TO order_user;

