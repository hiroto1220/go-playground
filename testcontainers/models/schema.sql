CREATE TABLE users (
  id CHAR(36) PRIMARY KEY,         
  name VARCHAR(100) NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO users (id, name, email) VALUES
('1', 'yamada taro', 'taro.yamada@example.com'),
('2', 'suzuki hanako', 'hanako.suzuki@example.com');