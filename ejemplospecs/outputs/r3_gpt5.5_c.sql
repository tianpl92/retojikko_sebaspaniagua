 CREATE TABLE customers (
        id BIGSERIAL PRIMARY KEY,
        email VARCHAR(255) NOT NULL,
        password_hash TEXT NOT NULL,
        full_name VARCHAR(255) NOT NULL,
        phone VARCHAR(32),
        default_shipping_address JSONB,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        deleted_at TIMESTAMPTZ
    );
    
    CREATE TABLE addresses (
        id BIGSERIAL PRIMARY KEY,
        customer_id BIGINT NOT NULL,
        recipient_name VARCHAR(255) NOT NULL,
        phone VARCHAR(32),
        address_line1 VARCHAR(255) NOT NULL,
        address_line2 VARCHAR(255),
        city VARCHAR(128) NOT NULL,
        state VARCHAR(128),
        postal_code VARCHAR(32) NOT NULL,
        country_code CHAR(2) NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        deleted_at TIMESTAMPTZ,
        CONSTRAINT addresses_customer_id_fk FOREIGN KEY (customer_id) REFERENCES customers (id) ON DELETE RESTRICT
    );
    
    CREATE TABLE products (
        id BIGSERIAL PRIMARY KEY,
        sku VARCHAR(128) NOT NULL,
        name VARCHAR(255) NOT NULL,
        description TEXT,
        price NUMERIC(12, 2) NOT NULL,
        stock_on_hand INTEGER NOT NULL DEFAULT 0,
        is_active BOOLEAN NOT NULL DEFAULT TRUE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        deleted_at TIMESTAMPTZ,
        CONSTRAINT products_sku_unique UNIQUE (sku),
        CONSTRAINT products_price_non_negative CHECK (price >= 0),
        CONSTRAINT products_stock_on_hand_non_negative CHECK (stock_on_hand >= 0)
    );
CREATE TABLE orders (
        id BIGSERIAL PRIMARY KEY,
        customer_id BIGINT NOT NULL,
        status VARCHAR(32) NOT NULL DEFAULT 'pending',
        total_amount NUMERIC(12, 2) NOT NULL DEFAULT 0,
        placed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        shipping_address_id BIGINT NOT NULL,
        billing_address_id BIGINT NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        deleted_at TIMESTAMPTZ,
        CONSTRAINT orders_customer_id_fk FOREIGN KEY (customer_id) REFERENCES customers (id) ON DELETE RESTRICT,
        CONSTRAINT orders_shipping_address_id_fk FOREIGN KEY (shipping_address_id) REFERENCES addresses (id) ON DELETE RESTRICT,
        CONSTRAINT orders_billing_address_id_fk FOREIGN KEY (billing_address_id) REFERENCES addresses (id) ON DELETE RESTRICT,
        CONSTRAINT orders_status_check CHECK (status IN ('pending', 'paid', 'shipped', 'delivered', 'cancelled')),
        CONSTRAINT orders_total_amount_non_negative CHECK (total_amount >= 0)
    );
    
    CREATE TABLE order_items (
        id BIGSERIAL PRIMARY KEY,
        order_id BIGINT NOT NULL,
        product_id BIGINT NOT NULL,
        quantity INTEGER NOT NULL,
        unit_price NUMERIC(12, 2) NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
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
        method VARCHAR(32) NOT NULL,
        status VARCHAR(32) NOT NULL DEFAULT 'initiated',
        processor_reference VARCHAR(255) NOT NULL,
        processed_at TIMESTAMPTZ,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        deleted_at TIMESTAMPTZ,
        CONSTRAINT payments_order_id_fk FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE RESTRICT,
        CONSTRAINT payments_processor_reference_unique UNIQUE (processor_reference),
        CONSTRAINT payments_method_check CHECK (method IN ('card', 'transfer', 'cash_on_delivery')),
        CONSTRAINT payments_status_check CHECK (status IN ('initiated', 'authorized', 'captured', 'failed', 'refunded')),
        CONSTRAINT payments_amount_valid_for_status CHECK (
            (status = 'refunded' AND amount < 0)
            OR (status <> 'refunded' AND amount >= 0)
        )
    );
 CREATE TABLE shipments (
        id BIGSERIAL PRIMARY KEY,
        order_id BIGINT NOT NULL,
        carrier VARCHAR(128) NOT NULL,
        tracking_number VARCHAR(255) NOT NULL,
        status VARCHAR(32) NOT NULL DEFAULT 'label_created',
        shipped_at TIMESTAMPTZ,
        delivered_at TIMESTAMPTZ,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        deleted_at TIMESTAMPTZ,
        CONSTRAINT shipments_order_id_fk FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE RESTRICT,
        CONSTRAINT shipments_tracking_number_unique UNIQUE (tracking_number),
        CONSTRAINT shipments_status_check CHECK (status IN ('label_created', 'in_transit', 'delivered', 'returned'))
    );
    
    CREATE TABLE shipment_items (
        id BIGSERIAL PRIMARY KEY,
        shipment_id BIGINT NOT NULL,
        order_item_id BIGINT NOT NULL,
        quantity INTEGER NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        deleted_at TIMESTAMPTZ,
        CONSTRAINT shipment_items_shipment_id_fk FOREIGN KEY (shipment_id) REFERENCES shipments (id) ON DELETE CASCADE,
        CONSTRAINT shipment_items_order_item_id_fk FOREIGN KEY (order_item_id) REFERENCES order_items (id) ON DELETE CASCADE,
        CONSTRAINT shipment_items_shipment_order_item_unique UNIQUE (shipment_id, order_item_id),
        CONSTRAINT shipment_items_quantity_positive CHECK (quantity > 0)
    );
    
 CREATE UNIQUE INDEX idx_customers_email_lower_unique ON customers (LOWER(email));
    CREATE INDEX idx_customers_deleted_at_active ON customers (id) WHERE deleted_at IS NULL;
    CREATE INDEX idx_addresses_customer_id ON addresses (customer_id);
    CREATE INDEX idx_addresses_deleted_at_active ON addresses (id) WHERE deleted_at IS NULL;
    CREATE INDEX idx_products_deleted_at_active ON products (id) WHERE deleted_at IS NULL;
    CREATE INDEX idx_orders_customer_id ON orders (customer_id);
    CREATE INDEX idx_orders_shipping_address_id ON orders (shipping_address_id);
    CREATE INDEX idx_orders_billing_address_id ON orders (billing_address_id);
    CREATE INDEX idx_orders_customer_id_placed_at_desc ON orders (customer_id, placed_at DESC);
    CREATE INDEX idx_orders_status_placed_at ON orders (status, placed_at);
    CREATE INDEX idx_orders_deleted_at_active ON orders (id) WHERE deleted_at IS NULL;
    CREATE INDEX idx_order_items_order_id ON order_items (order_id);
    CREATE INDEX idx_order_items_product_id ON order_items (product_id);
    CREATE INDEX idx_order_items_deleted_at_active ON order_items (id) WHERE deleted_at IS NULL;
    CREATE INDEX idx_payments_order_id ON payments (order_id);
    CREATE INDEX idx_payments_order_id_status ON payments (order_id, status);
    CREATE INDEX idx_payments_deleted_at_active ON payments (id) WHERE deleted_at IS NULL;
    CREATE INDEX idx_shipments_order_id ON shipments (order_id);
    CREATE INDEX idx_shipments_tracking_number ON shipments (tracking_number);
    CREATE INDEX idx_shipments_deleted_at_active ON shipments (id) WHERE deleted_at IS NULL;
    CREATE INDEX idx_shipment_items_shipment_id ON shipment_items (shipment_id);
    CREATE INDEX idx_shipment_items_order_item_id ON shipment_items (order_item_id);
    CREATE INDEX idx_shipment_items_deleted_at_active ON shipment_items (id) WHERE deleted_at IS NULL;

 COMMENT ON TABLE customers IS 'Account holders who place orders; customer deletion is represented by soft delete only.';
    COMMENT ON COLUMN customers.deleted_at IS 'Nullable soft-delete timestamp; when set, the customer is hidden from active workflows but historical orders remain intact.';
    
    COMMENT ON TABLE addresses IS 'Customer-owned shipping and billing addresses that can be reused across orders.';
    COMMENT ON COLUMN addresses.deleted_at IS 'Nullable soft-delete timestamp; when set, the address is no longer active but remains available for historical order references.';
    
    COMMENT ON TABLE products IS 'Sellable product catalog entries with current price and stock on hand.';
    COMMENT ON COLUMN products.deleted_at IS 'Nullable soft-delete timestamp; when set, the product is no longer active for sale but remains available for historical order items.';
    
    COMMENT ON TABLE orders IS 'Order headers containing customer, address, lifecycle status, and placement total snapshot.';
    COMMENT ON COLUMN orders.status IS 'Allowed values: pending, paid, shipped, delivered, cancelled. Lifecycle: pending -> paid -> shipped -> delivered; pending or paid may become cancelled.';
    COMMENT ON COLUMN orders.total_amount IS 'Snapshot of the sum of line items at placement time; do not recalculate historical totals from current product prices.';
    COMMENT ON COLUMN orders.deleted_at IS 'Nullable soft-delete timestamp; when set, the order is hidden from active workflows while preserving audit history.';
    
    COMMENT ON TABLE order_items IS 'Line items for orders, capturing product, quantity, and unit price at order placement.';
    COMMENT ON COLUMN order_items.unit_price IS 'Snapshot of product price at order placement time; never use products.price as a live join for historical order totals.';
    COMMENT ON COLUMN order_items.deleted_at IS 'Nullable soft-delete timestamp; when set, the line item is hidden from active workflows while preserving order history.';
    
    COMMENT ON TABLE payments IS 'Payment attempts, captures, failures, and refunds linked to orders; multiple rows per order are allowed.';
    COMMENT ON COLUMN payments.status IS 'Allowed values: initiated, authorized, captured, failed, refunded. Lifecycle: initiated -> authorized -> captured; initiated -> failed; captured -> refunded.';
    COMMENT ON COLUMN payments.deleted_at IS 'Nullable soft-delete timestamp; when set, the payment row is hidden from active workflows while preserving financial audit history.';
    
    COMMENT ON TABLE shipments IS 'Shipment parcels linked to orders; multiple shipments per order support split parcels.';
    COMMENT ON COLUMN shipments.status IS 'Allowed values: label_created, in_transit, delivered, returned. Lifecycle: label_created -> in_transit -> delivered; any state may become returned.';
    COMMENT ON COLUMN shipments.deleted_at IS 'Nullable soft-delete timestamp; when set, the shipment is hidden from active workflows while preserving fulfillment history.';
    
    COMMENT ON TABLE shipment_items IS 'Join table allocating quantities of order items to shipment parcels for partial shipments.';
    COMMENT ON COLUMN shipment_items.deleted_at IS 'Nullable soft-delete timestamp; when set, the shipment allocation is hidden from active workflows while preserving fulfillment history.';
