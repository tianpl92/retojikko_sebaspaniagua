    CREATE TABLE customers (
        customer_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
        email VARCHAR(255) NOT NULL UNIQUE,
        password_hash TEXT NOT NULL,
        first_name VARCHAR(100) NOT NULL,
        last_name VARCHAR(100) NOT NULL,
        phone VARCHAR(30),
        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    
    CREATE TABLE products (
        product_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
        sku VARCHAR(100) NOT NULL UNIQUE,
        name VARCHAR(255) NOT NULL,
        description TEXT,
        price NUMERIC(12, 2) NOT NULL,
        stock_quantity INTEGER NOT NULL DEFAULT 0,
        is_active BOOLEAN NOT NULL DEFAULT TRUE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        CONSTRAINT products_price_non_negative CHECK (price >= 0),
        CONSTRAINT products_stock_quantity_non_negative CHECK (stock_quantity >= 0)
    );    


    CREATE TABLE orders (
        order_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
        customer_id BIGINT NOT NULL,
        status VARCHAR(20) NOT NULL DEFAULT 'pending',
        currency CHAR(3) NOT NULL DEFAULT 'USD',
        subtotal_amount NUMERIC(12, 2) NOT NULL DEFAULT 0,
        tax_amount NUMERIC(12, 2) NOT NULL DEFAULT 0,
        shipping_amount NUMERIC(12, 2) NOT NULL DEFAULT 0,
        total_amount NUMERIC(12, 2) NOT NULL DEFAULT 0,
        placed_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        paid_at TIMESTAMPTZ,
        shipped_at TIMESTAMPTZ,
        delivered_at TIMESTAMPTZ,
        cancelled_at TIMESTAMPTZ,
        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        CONSTRAINT orders_customer_id_fk FOREIGN KEY (customer_id) REFERENCES customers (customer_id),
        CONSTRAINT orders_status_check CHECK (status IN ('pending', 'paid', 'shipped', 'delivered', 'cancelled')),
        CONSTRAINT orders_subtotal_amount_non_negative CHECK (subtotal_amount >= 0),
        CONSTRAINT orders_tax_amount_non_negative CHECK (tax_amount >= 0),
        CONSTRAINT orders_shipping_amount_non_negative CHECK (shipping_amount >= 0),
        CONSTRAINT orders_total_amount_non_negative CHECK (total_amount >= 0)
    );
    
    CREATE TABLE order_items (
        order_item_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
        order_id BIGINT NOT NULL,
        product_id BIGINT NOT NULL,
        quantity INTEGER NOT NULL,
        unit_price NUMERIC(12, 2) NOT NULL,
        line_total NUMERIC(12, 2) NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        CONSTRAINT order_items_order_id_fk FOREIGN KEY (order_id) REFERENCES orders (order_id) ON DELETE CASCADE,
        CONSTRAINT order_items_product_id_fk FOREIGN KEY (product_id) REFERENCES products (product_id),
        CONSTRAINT order_items_quantity_positive CHECK (quantity > 0),
        CONSTRAINT order_items_unit_price_non_negative CHECK (unit_price >= 0),
        CONSTRAINT order_items_line_total_non_negative CHECK (line_total >= 0)
    );
    

    CREATE TABLE payments (
        payment_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
        order_id BIGINT NOT NULL,
        status VARCHAR(20) NOT NULL DEFAULT 'pending',
        payment_method VARCHAR(50) NOT NULL,
        amount NUMERIC(12, 2) NOT NULL,
        currency CHAR(3) NOT NULL DEFAULT 'USD',
        provider VARCHAR(100),
        provider_transaction_id VARCHAR(255),
        paid_at TIMESTAMPTZ,
        failed_at TIMESTAMPTZ,
        refunded_at TIMESTAMPTZ,
        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        CONSTRAINT payments_order_id_fk FOREIGN KEY (order_id) REFERENCES orders (order_id) ON DELETE CASCADE,
        CONSTRAINT payments_status_check CHECK (status IN ('pending', 'authorized', 'paid', 'failed', 'refunded', 'cancelled')),
        CONSTRAINT payments_amount_positive CHECK (amount > 0),
        CONSTRAINT payments_provider_transaction_unique UNIQUE (provider, provider_transaction_id)
    );
    
    CREATE TABLE shipments (
        shipment_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
        order_id BIGINT NOT NULL,
        status VARCHAR(20) NOT NULL DEFAULT 'pending',
        carrier VARCHAR(100),
        tracking_number VARCHAR(255),
        shipped_at TIMESTAMPTZ,
        delivered_at TIMESTAMPTZ,
        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        CONSTRAINT shipments_order_id_fk FOREIGN KEY (order_id) REFERENCES orders (order_id) ON DELETE CASCADE,
        CONSTRAINT shipments_status_check CHECK (status IN ('pending', 'shipped', 'in_transit', 'delivered', 'failed', 'cancelled')),
        CONSTRAINT shipments_tracking_unique UNIQUE (carrier, tracking_number)
    );

