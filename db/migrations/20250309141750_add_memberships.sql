-- migrate:up
CREATE TABLE membership_plans (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL DEFAULT 0,
    billing_cycle VARCHAR(10) NOT NULL,
    features TEXT[],
    is_popular BOOLEAN NOT NULL DEFAULT false,
    max_folders INT DEFAULT NULL,
    custom_domain BOOLEAN NOT NULL DEFAULT false,
    advanced_analytics BOOLEAN NOT NULL DEFAULT false,
    stripe_product_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE memberships (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    membership_plan_id VARCHAR(255) NOT NULL REFERENCES membership_plans(id) ON DELETE CASCADE,
    status VARCHAR(32) NOT NULL,
    start_date TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    end_date TIMESTAMPTZ,
    stripe_subscription_id VARCHAR(255),
    stripe_subscription_interval VARCHAR(10),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE transactions (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    membership_id VARCHAR(255) NOT NULL REFERENCES memberships(id) ON DELETE CASCADE,
    amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    status VARCHAR(32) NOT NULL,
    stripe_payment_intent_id VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO membership_plans (id, name, price, billing_cycle, features, is_popular, max_folders, custom_domain, advanced_analytics, stripe_product_id)
VALUES 
    ('free', 'Free', 0, 'monthly', ARRAY['Basic photo uploads', 'Limited analytics'], false, 3, false, false, 'stripe_free_product_id'),
    ('monthly-unlimited', 'Unlimited', 19.99, 'monthly', ARRAY['Unlimited photo groups', 'Advanced analytics', 'Custom domain'], true, NULL, true, true, 'stripe_monthly_product_id'),
    ('yearly-unlimited', 'Unlimited', 199.99, 'yearly', ARRAY['Unlimited photo groups', 'Advanced analytics', 'Custom domain'], true, NULL, true, true, 'stripe_yearly_product_id');

-- migrate:down
DELETE FROM membership_plans WHERE id IN ('free', 'monthly-unlimited', 'yearly-unlimited');
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS memberships;
DROP TABLE IF EXISTS membership_plans;
