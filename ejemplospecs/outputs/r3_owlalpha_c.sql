-- =============================================================================
-- Order Management Schema — PostgreSQL 16
-- =============================================================================
-- This schema supports a small e-commerce backend with customers, products,
-- orders, payments, and shipments. It is designed for operational queries
-- (order history, revenue reports, low-stock alerts) and safe evolution
-- without destructive migrations.
--
-- Safe-change contract:
--   - ADD COLUMN must be nullable or DEFAULT a constant.
--   - No renames.
--   - No DROP COLUMN; deprecate via comments and stop writing.
--   - Future indexes must use CREATE INDEX CONCURRENTLY (not in this file).
-- =============================================================================

-- -----------------------------------------------------------------------------
-- customers
-- Account holder for the store. Email is unique (case-insensitive via CITEXT).
-- Soft-deletable: deleted_at marks deactivation without losing order history.
-- -----------------------------------------------------------------------------
CREATE TABLE customers (
    id              BIGSERIAL PRIMARY KEY,
    email           CITEXT NOT NULL UNIQUE,
    password_hash   TEXT NOT NULL,
    full_name       VARCHAR(255),
    phone           VARCHAR(50),
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at      TIMESTAMPTZ
);

COMMENT ON TABLE customers IS 'Account holders. Soft-deletable; orders and payments are preserved when a customer is deactivated.';
COMMENT ON COLUMN customers.email IS 'Unique login identifier. Case-insensitive via CITEXT.';
COMMENT ON COLUMN customers.password_hash IS 'Bcrypt/argon2 hash. Never store plaintext.';
COMMENT ON COLUMN customers.is_active IS 'FALSE disables login but preserves all historical data.';
COMMENT ON COLUMN customers.deleted_at IS 'Soft-delete timestamp. NULL means active. Set to now() to deactivate without losing referential integrity.';

-- -----------------------------------------------------------------------------
-- addresses
-- Customer addresses. A customer may have many; orders reference a snapshot via
-- shipping_address_id / billing_address_id to preserve history even if the
-- address is later edited or removed.
-- -----------------------------------------------------------------------------
CREATE TABLE addresses (
    id              BIGSERIAL PRIMARY KEY,
    customer_id     BIGINT NOT NULL,
    label           VARCHAR(100) DEFAULT 'home',
    street_line1    VARCHAR(255) NOT NULL,
    street_line2    VARCHAR(255),
    city            VARCHAR(100) NOT NULL,
    state           VARCHAR(100),
    postal_code     VARCHAR(20),
    country_code    CHAR(2) NOT NULL,
    is_default      BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at      TIMESTAMPTZ,
    CONSTRAINT fk_addresses_customer
        FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE RESTRICT
);

COMMENT ON TABLE addresses IS 'Customer addresses. Orders snapshot address_id to preserve historical billing/shipping locations.';
COMMENT ON COLUMN addresses.customer_id IS 'Owning customer. RESTRICT on delete — cannot remove a customer that has addresses.';
COMMENT ON COLUMN addresses.label IS 'Friendly name: home, office, etc.';
COMMENT ON COLUMN addresses.is_default IS 'TRUE for the customer''s default address (enforced in app, not SQL).';
COMMENT ON COLUMN addresses.deleted_at IS 'Soft-delete. NULL means active.';

-- -----------------------------------------------------------------------------
-- products
-- Product catalog. Price is the current live price; historical order prices are
-- captured in order_items.unit_price. Soft-deletable.
-- -----------------------------------------------------------------------------
CREATE TABLE products (
    id              BIGSERIAL PRIMARY KEY,
    sku             VARCHAR(100) NOT NULL UNIQUE,
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    price           NUMERIC(12, 2) NOT NULL CHECK (price >= 0),
    stock_on_hand   INTEGER NOT NULL DEFAULT 0 CHECK (stock_on_hand >= 0),
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at      TIMESTAMPTZ
);

COMMENT ON TABLE products IS 'Product catalog. Current price and stock. Historical prices are stored in order_items.unit_price, not joined from here.';
COMMENT ON COLUMN products.sku IS 'Stock-keeping unit. Unique across all products.';
COMMENT ON COLUMN products.price IS 'Current live price. Do NOT join to this for historical order totals — use order_items.unit_price instead.';
COMMENT ON COLUMN products.stock_on_hand IS 'Units available for sale. CHECK prevents negative inventory.';
COMMENT ON COLUMN products.is_active IS 'FALSE hides from catalog but preserves order history references.';
COMMENT ON COLUMN products.deleted_at IS 'Soft-delete. NULL means active.';

-- -----------------------------------------------------------------------------
-- orders
-- Order header. total_amount is a snapshot of the sum of line items at
-- placement time. shipping_address_id and billing_address_id point to
-- addresses at the time of order, preserving history even if the customer
-- later changes their address book.
-- -----------------------------------------------------------------------------
CREATE TABLE orders (
    id                      BIGSERIAL PRIMARY KEY,
    customer_id             BIGINT NOT NULL,
    status                  VARCHAR(32) NOT NULL DEFAULT 'pending'
                                CHECK (status IN ('pending', 'paid', 'shipped', 'delivered', 'cancelled')),
    total_amount            NUMERIC(12, 2) NOT NULL DEFAULT 0 CHECK (total_amount >= 0),
    placed_at               TIMESTAMPTZ NOT NULL DEFAULT now(),
    shipping_address_id     BIGINT,
    billing_address_id      BIGINT,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at              TIMESTAMPTZ,
    CONSTRAINT fk_orders_customer
        FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE RESTRICT,
    CONSTRAINT fk_orders_shipping_address
        FOREIGN KEY (shipping_address_id) REFERENCES addresses(id) ON DELETE RESTRICT,
    CONSTRAINT fk_orders_billing_address
        FOREIGN KEY (billing_address_id) REFERENCES addresses(id) ON DELETE RESTRICT
);

COMMENT ON TABLE orders IS 'Order headers. Immutable history — soft-delete only. Status transitions are enforced in application code.';
COMMENT ON COLUMN orders.customer_id IS 'Owning customer. RESTRICT on delete.';
COMMENT ON COLUMN orders.status IS 'Lifecycle: pending → paid → shipped → delivered. Terminal: pending → cancelled or paid → cancelled. Transitions enforced in app.';
COMMENT ON COLUMN orders.total_amount IS 'SNAPSHOT of SUM(order_items.subtotal) at placement time. Never recompute from live product prices.';
COMMENT ON COLUMN orders.placed_at IS 'When the customer submitted the order. Used for revenue reporting.';
COMMENT ON COLUMN orders.shipping_address_id IS 'Snapshot of the shipping address at order time. RESTRICT on delete to preserve history.';
COMMENT ON COLUMN orders.billing_address_id IS 'Snapshot of the billing address at order time. RESTRICT on delete to preserve history.';
COMMENT ON COLUMN orders.deleted_at IS 'Soft-delete. NULL means active.';

-- -----------------------------------------------------------------------------
-- order_items
-- Line items within an order. unit_price is a snapshot from products.price at
-- the moment the order was placed. subtotal = quantity * unit_price.
-- -----------------------------------------------------------------------------
CREATE TABLE order_items (
    id              BIGSERIAL PRIMARY KEY,
    order_id        BIGINT NOT NULL,
    product_id      BIGINT NOT NULL,
    quantity        INTEGER NOT NULL CHECK (quantity > 0),
    unit_price      NUMERIC(12, 2) NOT NULL CHECK (unit_price >= 0),
    subtotal        NUMERIC(12, 2) NOT NULL CHECK (subtotal >= 0),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at      TIMESTAMPTZ,
    CONSTRAINT fk_order_items_order
        FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE RESTRICT,
    CONSTRAINT fk_order_items_product
        FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE RESTRICT
);

COMMENT ON TABLE order_items IS 'Line items of an order. unit_price is a snapshot; never join back to products.price for historical totals.';
COMMENT ON COLUMN order_items.order_id IS 'Parent order. RESTRICT on delete — order history is immutable.';
COMMENT ON COLUMN order_items.product_id IS 'Purchased product. RESTRICT on delete — cannot remove a product that has been sold.';
COMMENT ON COLUMN order_items.quantity IS 'Units ordered. Must be positive.';
COMMENT ON COLUMN order_items.unit_price IS 'SNAPSHOT of products.price at order time. Do NOT join to products.price for revenue calculations.';
COMMENT ON COLUMN order_items.subtotal IS 'quantity * unit_price. Stored as a snapshot for fast reporting.';
COMMENT ON COLUMN order_items.deleted_at IS 'Soft-delete. NULL means active.';

-- -----------------------------------------------------------------------------
-- payments
-- Payment attempts for an order. Multiple rows per order are expected (retries,
-- partial captures). A refund is a new row with negative amount and status
-- 'refunded'. The order is considered paid when at least one payment row has
-- status = 'captured'.
-- -----------------------------------------------------------------------------
CREATE TABLE payments (
    id                      BIGSERIAL PRIMARY KEY,
    order_id                BIGINT NOT NULL,
    amount                  NUMERIC(12, 2) NOT NULL CHECK (amount >= 0),
    method                  VARCHAR(32) NOT NULL
                                CHECK (method IN ('card', 'transfer', 'cash_on_delivery')),
    status                  VARCHAR(32) NOT NULL DEFAULT 'initiated'
                                CHECK (status IN ('initiated', 'authorized', 'captured', 'failed', 'refunded')),
    processor_reference     VARCHAR(255) UNIQUE,
    processed_at            TIMESTAMPTZ,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at              TIMESTAMPTZ,
    CONSTRAINT fk_payments_order
        FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE RESTRICT
);

COMMENT ON TABLE payments IS 'Payment attempts for an order. Multiple rows allowed (retries, refunds). Order is paid when any row has status = captured.';
COMMENT ON COLUMN payments.order_id IS 'Parent order. RESTRICT on delete.';
COMMENT ON COLUMN payments.amount IS 'Payment amount. Must be >= 0. Refunds are separate rows with status = refunded.';
COMMENT ON COLUMN payments.method IS 'Payment method: card, transfer, or cash_on_delivery.';
COMMENT ON COLUMN payments.status IS 'Lifecycle: initiated → authorized → captured, or initiated → failed, or captured → refunded. Transitions enforced in app.';
COMMENT ON COLUMN payments.processor_reference IS 'Unique reference from the payment processor. Used for idempotency and reconciliation.';
COMMENT ON COLUMN payments.processed_at IS 'When the payment was finalized by the processor. NULL for pending/failed attempts.';
COMMENT ON COLUMN payments.deleted_at IS 'Soft-delete. NULL means active.';

-- -----------------------------------------------------------------------------
-- shipments
-- Shipment records for an order. An order may have multiple shipments (split
-- parcels). Tracking number is unique across all shipments.
-- -----------------------------------------------------------------------------
CREATE TABLE shipments (
    id                      BIGSERIAL PRIMARY KEY,
    order_id                BIGINT NOT NULL,
    carrier                 VARCHAR(100),
    tracking_number         VARCHAR(255) UNIQUE,
    status                  VARCHAR(32) NOT NULL DEFAULT 'label_created'
                                CHECK (status IN ('label_created', 'in_transit', 'delivered', 'returned')),
    shipped_at              TIMESTAMPTZ,
    delivered_at            TIMESTAMPTZ,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at              TIMESTAMPTZ,
    CONSTRAINT fk_shipments_order
        FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE RESTRICT
);

COMMENT ON TABLE shipments IS 'Shipment parcels for an order. Multiple shipments per order are supported (split shipping).';
COMMENT ON COLUMN shipments.order_id IS 'Parent order. RESTRICT on delete.';
COMMENT ON COLUMN shipments.carrier IS 'Shipping carrier name (e.g. UPS, FedEx, DHL).';
COMMENT ON COLUMN shipments.tracking_number IS 'Unique tracking number. Used for customer-facing tracking lookups.';
COMMENT ON COLUMN shipments.status IS 'Lifecycle: label_created → in_transit → delivered, or any → returned. Transitions enforced in app.';
COMMENT ON COLUMN shipments.shipped_at IS 'When the parcel was handed to the carrier.';
COMMENT ON COLUMN shipments.delivered_at IS 'When the parcel was confirmed delivered. NULL until delivery.';
COMMENT ON COLUMN shipments.deleted_at IS 'Soft-delete. NULL means active.';

-- -----------------------------------------------------------------------------
-- shipment_items
-- Join table between shipments and order_items. Supports partial shipments:
-- a single order_item (qty 10) can be split across shipments (6 + 4).
-- The quantity column on each row indicates how many units of that order_item
-- are in this shipment.
-- -----------------------------------------------------------------------------
CREATE TABLE shipment_items (
    id              BIGSERIAL PRIMARY KEY,
    shipment_id     BIGINT NOT NULL,
    order_item_id   BIGINT NOT NULL,
    quantity        INTEGER NOT NULL CHECK (quantity > 0),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at      TIMESTAMPTZ,
    CONSTRAINT fk_shipment_items_shipment
        FOREIGN KEY (shipment_id) REFERENCES shipments(id) ON DELETE CASCADE,
    CONSTRAINT fk_shipment_items_order_item
        FOREIGN KEY (order_item_id) REFERENCES order_items(id) ON DELETE CASCADE,
    CONSTRAINT uq_shipment_items_pair
        UNIQUE (shipment_id, order_item_id)
);

COMMENT ON TABLE shipment_items IS 'Many-to-many join between shipments and order_items. Supports splitting a single order_item across multiple parcels.';
COMMENT ON COLUMN shipment_items.shipment_id IS 'Parent shipment. CASCADE on delete — this join row has no independent meaning.';
COMMENT ON COLUMN shipment_items.order_item_id IS 'The order_item being shipped in this parcel. CASCADE on delete.';
COMMENT ON COLUMN shipment_items.quantity IS 'How many units of this order_item are in this shipment. Must be positive.';
COMMENT ON COLUMN shipment_items.deleted_at IS 'Soft-delete. NULL means active.';

-- =============================================================================
-- INDEXES
-- =============================================================================

-- Partial index for active customers (excludes soft-deleted rows)
CREATE INDEX idx_customers_active ON customers (email) WHERE deleted_at IS NULL;

CREATE INDEX idx_addresses_customer_id ON addresses (customer_id);
CREATE INDEX idx_addresses_active ON addresses (customer_id) WHERE deleted_at IS NULL;

-- Partial index for active products
CREATE INDEX idx_products_sku ON products (sku) WHERE deleted_at IS NULL;
CREATE INDEX idx_products_active ON products (id) WHERE deleted_at IS NULL;

-- Composite index: customer order history (most recent first)
CREATE INDEX idx_orders_customer_placed ON orders (customer_id, placed_at DESC);

-- Composite index: operational dashboards filtering by status and date
CREATE INDEX idx_orders_status_placed ON orders (status, placed_at);

-- Partial index for active orders
CREATE INDEX idx_orders_active ON orders (id) WHERE deleted_at IS NULL;

CREATE INDEX idx_order_items_order_id ON order_items (order_id);
CREATE INDEX idx_order_items_product_id ON order_items (product_id);
CREATE INDEX idx_order_items_active ON order_items (order_id) WHERE deleted_at IS NULL;

-- Composite index: payment lookups by order and status
CREATE INDEX idx_payments_order_status ON payments (order_id, status);

-- Partial index for active payments
CREATE INDEX idx_payments_active ON payments (order_id) WHERE deleted_at IS NULL;

CREATE INDEX idx_shipments_order_id ON shipments (order_id);

-- Hot-path lookup: tracking number search
CREATE INDEX idx_shipments_tracking ON shipments (tracking_number);

-- Partial index for active shipments
CREATE INDEX idx_shipments_active ON shipments (order_id) WHERE deleted_at IS NULL;

CREATE INDEX idx_shipment_items_shipment_id ON shipment_items (shipment_id);
CREATE INDEX idx_shipment_items_order_item_id ON shipment_items (order_item_id);
