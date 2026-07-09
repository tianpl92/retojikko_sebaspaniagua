CREATE TABLE customers (
        id BIGSERIAL PRIMARY KEY,
        email VARCHAR(255) NOT NULL,
        password_hash TEXT NOT NULL,
        first_name VARCHAR(100) NOT NULL,
        last_name VARCHAR(100) NOT NULL,
        address_line1 VARCHAR(255) NOT NULL,
        address_line2 VARCHAR(255),
        city VARCHAR(100) NOT NULL,
        state VARCHAR(100),
        postal_code VARCHAR(30) NOT NULL,
        country_code CHAR(2) NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        deleted_at TIMESTAMPTZ,
        CONSTRAINT customers_email_unique UNIQUE (email)
    );
    
    CREATE TABLE products (
        id BIGSERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        sku VARCHAR(100) NOT NULL,
        price NUMERIC(12, 2) NOT NULL,
        stock_on_hand INTEGER NOT NULL DEFAULT 0,
        is_active BOOLEAN NOT NULL DEFAULT TRUE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        deleted_at TIMESTAMPTZ,
        CONSTRAINT products_sku_unique UNIQUE (sku),
        CONSTRAINT products_price_non_negative CHECK (price >= 0),
        CONSTRAINT products_stock_on_hand_non_negative CHECK (stock_on_hand >= 0)
    );

 CREATE TABLE orders (
        id BIGSERIAL PRIMARY KEY,
        customer_id BIGINT NOT NULL,
        status VARCHAR(20) NOT NULL DEFAULT 'pending',
        total_amount NUMERIC(12, 2) NOT NULL DEFAULT 0,
        placed_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        paid_at TIMESTAMPTZ,
        shipped_at TIMESTAMPTZ,
        delivered_at TIMESTAMPTZ,
        cancelled_at TIMESTAMPTZ,
        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        deleted_at TIMESTAMPTZ,
        CONSTRAINT orders_customer_id_fk FOREIGN KEY (customer_id) REFERENCES customers (id) ON DELETE RESTRICT,
        CONSTRAINT orders_status_check CHECK (status IN ('pending', 'paid', 'shipped', 'delivered', 'cancelled')),
        CONSTRAINT orders_total_amount_non_negative CHECK (total_amount >= 0)
    );
    
    CREATE TABLE order_items (
        id BIGSERIAL PRIMARY KEY,
        order_id BIGINT NOT NULL,
        product_id BIGINT NOT NULL,
        quantity INTEGER NOT NULL,
        unit_price NUMERIC(12, 2) NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        deleted_at TIMESTAMPTZ,
        CONSTRAINT order_items_order_id_fk FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE RESTRICT,
        CONSTRAINT order_items_product_id_fk FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE RESTRICT,
        CONSTRAINT order_items_quantity_positive CHECK (quantity > 0),
        CONSTRAINT order_items_unit_price_non_negative CHECK (unit_price >= 0)
    );

   CREATE TABLE payments (
        id BIGSERIAL PRIMARY KEY,
        order_id BIGINT NOT NULL,
        amount NUMERIC(12, 2) NOT NULL,
        method VARCHAR(50) NOT NULL,
        status VARCHAR(20) NOT NULL DEFAULT 'pending',
        processor_reference VARCHAR(255) NOT NULL,
        processed_at TIMESTAMPTZ,
        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        deleted_at TIMESTAMPTZ,
        CONSTRAINT payments_order_id_fk FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE RESTRICT,
        CONSTRAINT payments_processor_reference_unique UNIQUE (processor_reference),
        CONSTRAINT payments_status_check CHECK (status IN ('pending', 'authorized', 'paid', 'failed', 'cancelled')),
        CONSTRAINT payments_amount_positive CHECK (amount > 0)
    );
    
    CREATE TABLE shipments (
        id BIGSERIAL PRIMARY KEY,
        order_id BIGINT NOT NULL,
        carrier VARCHAR(100) NOT NULL,
        tracking_number VARCHAR(255) NOT NULL,
        status VARCHAR(20) NOT NULL DEFAULT 'pending',
        shipped_at TIMESTAMPTZ,
        delivered_at TIMESTAMPTZ,
        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        deleted_at TIMESTAMPTZ,
        CONSTRAINT shipments_order_id_fk FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE RESTRICT,
        CONSTRAINT shipments_status_check CHECK (status IN ('pending', 'shipped', 'in_transit', 'delivered', 'failed', 'cancelled'))
    );
    
    CREATE INDEX idx_orders_customer_id ON orders (customer_id);
    CREATE INDEX idx_orders_status ON orders (status);
    CREATE INDEX idx_order_items_order_id ON order_items (order_id);
    CREATE INDEX idx_order_items_product_id ON order_items (product_id);
    CREATE INDEX idx_payments_order_id ON payments (order_id);
    CREATE INDEX idx_shipments_order_id ON shipments (order_id);
    CREATE INDEX idx_shipments_tracking_number ON shipments (tracking_number);

