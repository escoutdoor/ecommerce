CREATE TYPE roles as ENUM('admin', 'customer');

CREATE TABLE IF NOT EXISTS users(
    "id" SERIAL PRIMARY KEY,
    "email" VARCHAR UNIQUE NOT NULL,
    "first_name" VARCHAR NOT NULL,
    "last_name" VARCHAR NULL,
    "date_of_birth" DATE NULL,
    "password" TEXT NOT NULL,
    "role" roles NOT NULL DEFAULT 'customer',
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS categories(
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS products(
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR NOT NULL,
    "description" TEXT,
    "price" DECIMAL(10, 2) NOT NULL,
    "category_id" INTEGER,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY ("category_id") REFERENCES categories ("id") ON DELETE SET NULL
);

CREATE TYPE order_status as ENUM('pending', 'processing', 'shipped', 'delivered', 'cancelled');

CREATE TABLE IF NOT EXISTS orders (
    "id" SERIAL PRIMARY KEY,
    "total" DECIMAL(10, 2) NOT NULL,
    "user_id" INTEGER NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT "fk_user" FOREIGN KEY ("user_id")
        REFERENCES users ("id")
        ON DELETE SET NULL
);

CREATE TABLE shipping_details (
    "id" SERIAL PRIMARY KEY,
    "address_line1" VARCHAR NOT NULL,
    "address_line2" VARCHAR,
    "postal_code" VARCHAR(20),
    "city" VARCHAR(100) NOT NULL,
    "country" VARCHAR(100) NOT NULL,
    "notes" TEXT DEFAULT '',
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS order_items (
    "id" SERIAL PRIMARY KEY,
    "status" order_status NOT NULL DEFAULT 'pending',
    "product_id" INTEGER,
    "order_id" INTEGER NOT NULL,
    "shipping_details_id" INTEGER NOT NULL,
    "quantity" INTEGER NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT "fk_order" FOREIGN KEY ("order_id")
        REFERENCES orders ("id")
        ON DELETE CASCADE,
    CONSTRAINT "fk_product" FOREIGN KEY ("product_id")
        REFERENCES products ("id")
        ON DELETE SET NULL,
    CONSTRAINT "fk_shipping_details" FOREIGN KEY ("shipping_details_id")
        REFERENCES shipping_details ("id")
        ON DELETE RESTRICT
);

