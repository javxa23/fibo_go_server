-- CHECKED
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    points INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    category_id INT,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    is_approved BOOLEAN DEFAULT FALSE,
    view_count INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_category FOREIGN KEY (category_id) REFERENCES categories(id)
);

-- NOT CHECKED

-- CREATE TABLE comments (
--     id SERIAL PRIMARY KEY,
--     post_id INT NOT NULL,
--     user_id INT,
--     parent_comment_id INT DEFAULT NULL,
--     content TEXT NOT NULL,
--     created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
--     CONSTRAINT fk_post_comment FOREIGN KEY (post_id) REFERENCES posts(id),
--     CONSTRAINT fk_user_comment FOREIGN KEY (user_id) REFERENCES users(id),
--     CONSTRAINT fk_parent_comment FOREIGN KEY (parent_comment_id) REFERENCES comments(id)
-- );

-- CREATE TABLE likes (
--     id SERIAL PRIMARY KEY,
--     post_id INT NOT NULL,
--     user_id INT,
--     created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
--     CONSTRAINT fk_post_like FOREIGN KEY (post_id) REFERENCES posts(id),
--     CONSTRAINT fk_user_like FOREIGN KEY (user_id) REFERENCES users(id)
-- );

-- CREATE TABLE salaries (
--     id SERIAL PRIMARY KEY,
--     user_id INT NOT NULL,
--     month_year DATE NOT NULL,
--     reputation_points INT NOT NULL,
--     salary_amount NUMERIC(10, 2) NOT NULL,
--     created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
--     CONSTRAINT fk_user_salary FOREIGN KEY (user_id) REFERENCES users(id)
-- );
