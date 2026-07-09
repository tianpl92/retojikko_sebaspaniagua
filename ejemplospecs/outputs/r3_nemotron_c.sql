CREATE EXTENSION IF NOT EXISTS "citex";
    
    CREATE TABLE customers (
        id BIGSERIAL PRIMARY KEY,
        email CITEXT NOT NULL UNIQUE,
        hashed_password VARCHAR(255) NOT NULL,
        full_name VARCHAR(255) NOT NULL,
        phone VARCHAR(20),
        default_shipping_address JSONB,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        deleted_at TIMESTAMPTZ NULL
    );
    
    COMMENT ON TABLE customers IS 'Account holder.';
    COMMENT ON COLUMN customers.email IS 'Unique email address (case-insensitive).';
    COMMENT ON COLUMN customers.hashed_password IS 'Hashed password for authentication.';
    COMMENT ON COLUMN customers.full_name IS 'Customer''s full name.';
    COMMENT ON COLUMN customers.phone IS 'Phone number.';
    COMMENT ON COLUMN customers.default_shipping_address IS 'Default shipping address as JSON denormalised for convenience.';
    COMMENT ON COLUMN customers.created_at IS 'Timestamp when the customer was created.';
    COMMENT ON COLUMN customers.updated_at IS 'Timestamp when the customer was last updated.';
    COMMENT ON COLUMN customers.deleted_at IS 'Timestamp when the customer was soft deleted (nullable).';
    
    CREATE TABLE addresses (
        id BIGSERIAL PRIMARY KEY,
        customer_id BIGINT NOT NULL REFERENCES customers(id) ON DELETE RESTRICT,
        label VARCHAR(50),
        street_address VARCHAR(255) NOT NULL,
        city VARCHAR(100) NOT NULL,
        state VARCHAR(100),
        postal_code VARCHAR(20),
        country VARCHAR(100) NOT NULL,
        is_default_shipping BOOLEAN DEFAULT false,
        is_default_billing BOOLEAN DEFAULT false,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        deleted_at TIMESTAMPTZ NULL
    );

COMMENT ON TABLE addresses IS 'Customer addresses (multiple per customer, reusable across orders).';
    COMMENT ON COLUMN addresses.customer_id IS 'Foreign key to customer.';
    COMMENT ON COLUMN addresses.label IS 'Optional label (e.g., ''Home'', ''Work'').';
    COMMENT ON COLUMN addresses.street_address IS 'Street address.';
    COMMENT ON COLUMN addresses.city IS 'City.';
    COMMENT ON COLUMN addresses.state IS 'State or province.';
    COMMENT ON COLUMN addresses.postal_code IS 'Postal code.';
    COMMENT ON COLUMN addresses.country IS 'Country.';
    COMMENT ON COLUMN addresses.is_default_shipping IS 'Boolean flag if this is the customer''s default shipping address.';
    COMMENT ON COLUMN addresses.is_default_billing IS 'Boolean flag if this is the customer''s default billing address.';
    COMMENT ON COLUMN addresses.created_at IS 'Timestamp when the address was created.';
    COMMENT ON COLUMN addresses.updated_at IS 'Timestamp when the address was last updated.';
    COMMENT ON COLUMN addresses.deleted_at IS 'Timestamp when the address was soft deleted (nullable).';
    
    CREATE TABLE products (
        id BIGSERIAL PRIMARY KEY,
        sku VARCHAR(100) NOT NULL UNIQUE,
        name VARCHAR(255) NOT NULL,
        description TEXT,
        price NUMERIC(12,2) NOT NULL CHECK (price >= 0),
        stock_on_hand INTEGER NOT NULL DEFAULT 0 CHECK (stock_on_hand >= 0),
        created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        deleted_at TIMESTAMPTZ NULL
    );
    
    COMMENT ON TABLE products IS 'Product catalog.';
    COMMENT ON COLUMN products.sku IS 'Unique stock keeping unit.';
    COMMENT ON COLUMN products.name IS 'Product name.';
    COMMENT ON COLUMN products.description IS 'Product description.';
    COMMENT ON COLUMN products.price IS 'Current price (snapshot in order_items uses historical price).';
    COMMENT ON COLUMN products.stock_on_hand IS 'Current stock quantity.';
    COMMENT ON COLUMN products.created_at IS 'Timestamp when the product was created.';
    COMMENT ON COLUMN products.updated_at IS 'Timestamp when the product was last updated.';
    COMMENT ON COLUMN products.deleted_at IS 'Timestamp when the product was soft deleted (nullable).';


   CREATE TABLE orders (
        id BIGSERIAL PRIMARY KEY,
        customer_id BIGINT NOT NULL REFERENCES customers(id) ON DELETE RESTRICT,
        status VARCHAR(32) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'paid', 'shipped', 'delivered', 'cancelled')),
        total_amount NUMERIC(12,2) NOT NULL CHECK (total_amount >= 0),
        placed_at TIMESTAMPTZ NOT NULL,
        shipping_address_id BIGINT NOT NULL REFERENCES addresses(id) ON DELETE RESTRICT,
        billing_address_id BIGINT NOT NULL REFERENCES addresses(id) ON DELETE RESTRICT,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        deleted_at TIMESTAMPTZ NULL
    );
    
    COMMENT ON TABLE orders IS 'Order header.';
    COMMENT ON COLUMN orders.customer_id IS 'Foreign key to customer who placed the order.';
    COMMENT ON COLUMN orders.status IS 'Order lifecycle: pending → paid → shipped → delivered, or cancelled from pending/paid.';
    COMMENT ON COLUMN orders.total_amount IS 'Snapshot of order total at placement (sum of line items). Not a live join.';
    COMMENT ON COLUMN orders.placed_at IS 'Timestamp when the order was placed.';
    COMMENT ON COLUMN orders.shipping_address_id IS 'Foreign key to address used for shipping.';
    COMMENT ON COLUMN orders.billing_address_id IS 'Foreign key to address used for billing.';
    COMMENT ON COLUMN orders.created_at IS 'Timestamp when the order record was created.';
    COMMENT ON COLUMN orders.updated_at IS 'Timestamp when the order record was last updated.';
    COMMENT ON COLUMN orders.deleted_at IS 'Timestamp when the order was soft deleted (nullable).';
    
    CREATE TABLE order_items (
        id BIGSERIAL PRIMARY KEY,
        order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE RESTRICT,
        product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
        quantity INTEGER NOT NULL CHECK (quantity > 0),
        unit_price NUMERIC(12,2) NOT NULL CHECK (unit_price >= 0),
        created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        deleted_at TIMESTAMPTZ NULL
    );
    
    COMMENT ON TABLE order_items IS 'Line items for an order.';
    COMMENT ON COLUMN order_items.order_id IS 'Foreign key to parent order.';
    COMMENT ON COLUMN order_items.product_id IS 'Foreign key to product (snapshot of product at order time).';
    COMMENT ON COLUMN order_items.quantity IS 'Quantity of product ordered.';
    COMMENT ON COLUMN order_items.unit_price IS 'Unit price at order time (snapshot, not live join to products.price).';
    COMMENT ON COLUMN order_items.created_at IS 'Timestamp when the line item was created.';
    COMMENT ON COLUMN order_items.updated_at IS 'Timestamp when the line item was last updated.';
    COMMENT ON COLUMN order_items.deleted_at IS 'Timestamp when the line item was soft deleted (nullable).';


  CREATE TABLE payments (
        id BIGSERIAL PRIMARY KEY,
        order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE RESTRICT,
        amount NUMERIC(12,2) NOT NULL CHECK (amount >= 0),
        method VARCHAR(32),
        processor_reference VARCHAR(255) UNIQUE,
        status VARCHAR(32) NOT NULL DEFAULT 'initiated' CHECK (status IN ('initiated', 'authorized', 'captured', 'failed', 'refunded')),
        processed_at TIMESTAMPTZ NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        deleted_at TIMESTAMPTZ NULL
    );
    
    COMMENT ON TABLE payments IS 'Payment records linked to orders.';
    COMMENT ON COLUMN payments.order_id IS 'Foreign key to parent order.';
    COMMENT ON COLUMN payments.amount IS 'Payment amount (positive for charges, negative for refunds).';
    COMMENT ON COLUMN payments.method IS 'Payment method (e.g., card, transfer, cash_on_delivery).';
    COMMENT ON COLUMN payments.processor_reference IS 'Unique reference from payment processor.';
    COMMENT ON COLUMN payments.status IS 'Payment lifecycle: initiated → authorized → captured, or initiated → failed, or captured → refunded.';
    COMMENT ON COLUMN payments.processed_at IS 'Timestamp when payment was processed by processor.';
    COMMENT ON COLUMN payments.created_at IS 'Timestamp when the payment record was created.';
    COMMENT ON COLUMN payments.updated_at IS 'Timestamp when the payment record was last updated.';
    COMMENT ON COLUMN payments.deleted_at IS 'Timestamp when the payment was soft deleted (nullable).';
    
    CREATE TABLE shipments (
        id BIGSERIAL PRIMARY KEY,
        order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE RESTRICT,
        carrier VARCHAR(100),
        tracking_number VARCHAR(100) UNIQUE,
        status VARCHAR(32) NOT NULL DEFAULT 'label_created' CHECK (status IN ('label_created', 'in_transit', 'delivered', 'returned')),
        shipped_at TIMESTAMPTZ NULL,
        delivered_at TIMESTAMPTZ NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        deleted_at TIMESTAMPTZ NULL
    );
    
    COMMENT ON TABLE shipments IS 'Shipment records linked to orders (multiple parcels per order).';
    COMMENT ON COLUMN shipments.order_id IS 'Foreign key to parent order.';
    COMMENT ON COLUMN shipments.carrier IS 'Shipping carrier (e.g., UPS, FedEx).';
    COMMENT ON COLUMN shipments.tracking_number IS 'Tracking number from carrier.';
    COMMENT ON COLUMN shipments.status IS 'Shipment lifecycle: label_created → in_transit → delivered, or any state → returned.';
    COMMENT ON COLUMN shipments.shipped_at IS 'Timestamp when shipment was shipped.';
    COMMENT ON COLUMN shipments.delivered_at IS 'Timestamp when shipment was delivered.';
    COMMENT ON COLUMN shipments.created_at IS 'Timestamp when the shipment record was created.';
    COMMENT ON COLUMN shipments.updated_at IS 'Timestamp when the shipment record was last updated.';
    COMMENT ON COLUMN shipments.deleted_at IS 'Timestamp when the shipment was soft deleted (nullable).';


   CREATE TABLE shipment_items (
        id BIGSERIAL PRIMARY KEY,
        shipment_id BIGINT NOT NULL REFERENCES shipments(id) ON DELETE CASCADE,
        order_item_id BIGINT NOT NULL REFERENCES order_items(id) ON DELETE CASCADE,
        quantity INTEGER NOT NULL CHECK (quantity > 0),
        created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        deleted_at TIMESTAMPTZ NULL
    );
    
    COMMENT ON TABLE shipment_items IS 'Many-to-many between shipments and order_items (supports split parcels).';
    COMMENT ON COLUMN shipment_items.shipment_id IS 'Foreign key to shipment.';
    COMMENT ON COLUMN shipment_items.order_item_id IS 'Foreign key to order item.';
    COMMENT ON COLUMN shipment_items.quantity IS 'Quantity of the order item included in this shipment.';
    COMMENT ON COLUMN shipment_items.created_at IS 'Timestamp when the shipment item was created.';
    COMMENT ON COLUMN shipment_items.updated_at IS 'Timestamp when the shipment item was last updated.';
    COMMENT ON COLUMN shipment_items.deleted_at IS 'Timestamp when the shipment item was soft deleted (nullable).';
    
    CREATE INDEX idx_customers_email ON customers(email);
    CREATE INDEX idx_addresses_customer_id ON addresses(customer_id);
    CREATE INDEX idx_products_sku ON products(sku);
    CREATE INDEX idx_orders_customer_id ON orders(customer_id);
    CREATE INDEX idx_orders_status ON orders(status);
    CREATE INDEX idx_orders_customer_placed ON orders(customer_id, placed_at DESC);
    CREATE INDEX idx_orders_status_placed ON orders(status, placed_at);
    CREATE INDEX idx_order_items_order_id ON order_items(order_id);
    CREATE INDEX idx_order_items_product_id ON order_items(product_id);
    CREATE INDEX idx_payments_order_id ON payments(order_id);
    CREATE INDEX idx_payments_status ON payments(status);
    CREATE INDEX idx_payments_order_status ON payments(order_id, status);
    CREATE INDEX idx_shipments_order_id ON shipments(order_id);
    CREATE INDEX idx_shipments_status ON shipments(status);
    CREATE INDEX idx_shipments_tracking_number ON shipments(tracking_number);
    CREATE INDEX idx_shipment_items_shipment_id ON shipment_items(shipment_id);
    CREATE INDEX idx_shipment_items_order_item_id ON shipment_items(order_item_id);
    CREATE INDEX idx_shipment_items_unique ON shipment_items(shipment_id, order_item_id);
    
    CREATE INDEX idx_customers_deleted_at ON customers(deleted_at) WHERE deleted_at IS NULL;
    CREATE INDEX idx_addresses_deleted_at ON addresses(deleted_at) WHERE deleted_at IS NULL;
    CREATE INDEX idx_products_deleted_at ON products(deleted_at) WHERE deleted_at IS NULL;
    CREATE INDEX idx_orders_deleted_at ON orders(deleted_at) WHERE deleted_at IS NULL;
    CREATE INDEX idx_order_items_deleted_at ON order_items(deleted_at) WHERE deleted_at IS NULL;
    CREATE INDEX idx_payments_deleted_at ON payments(deleted_at) WHERE deleted_at IS NULL;
    CREATE INDEX idx_shipments_deleted_at ON shipments(deleted_at) WHERE deleted_at IS NULL;
    CREATE INDEX idx_shipment_items_deleted_at ON shipment_items(deleted_at) WHERE deleted_at IS NULL;

