-- migrate:up
CREATE TABLE portfolios (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    theme VARCHAR(255) NOT NULL DEFAULT 'light',
    description TEXT,
    columns INT NOT NULL DEFAULT 3,
    gap INT NOT NULL DEFAULT 16,
    rounded_corners BOOLEAN NOT NULL DEFAULT true,
    show_captions BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE profiles (
    id VARCHAR(255) PRIMARY KEY,
    portfolio_id VARCHAR(255) REFERENCES portfolios(id) ON DELETE CASCADE,
    name VARCHAR(150) NOT NULL,
    title VARCHAR(150),
    bio TEXT,
    email VARCHAR(150),
    instagram VARCHAR(150),
    website VARCHAR(255)
);

CREATE TABLE folders (
    id VARCHAR(255) PRIMARY KEY,
    portfolio_id VARCHAR(255) REFERENCES portfolios(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    cover_id INT,
    columns INT NOT NULL DEFAULT 3,
    gap INT NOT NULL DEFAULT 16,
    rounded_corners BOOLEAN NOT NULL DEFAULT true,
    show_captions BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE photos (
    id VARCHAR(255) PRIMARY KEY,
    folder_id VARCHAR(255),
    src TEXT NOT NULL,
    alt TEXT,
    caption TEXT,
    public_id VARCHAR(255),
    sort_index INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- migrate:down
DROP TABLE IF EXISTS photos;
DROP TABLE IF EXISTS folders;
DROP TABLE IF EXISTS profiles;
DROP TABLE IF EXISTS portfolios;
