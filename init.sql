CREATE SCHEMA users;
CREATE TABLE users.users (
  id SERIAL PRIMARY KEY NOT NULL,
  created_at TIMESTAMP DEFAULT NOW() NOT NULL,
  updated_at TIMESTAMP DEFAULT NOW() NOT NULL,
  google_id VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL,
  name VARCHAR(255) NOT NULL,
  picture VARCHAR(255) NOT NULL
);
INSERT INTO users.users (google_id, email, name, picture) VALUES
  ('123456789012345678901', 'john.doe@example.com', 'John Doe',
  'https://lh3.googleusercontent.com/a/default-user=s96-c'),
  ('234567890123456789012', 'jane.smith@example.com', 'Jane Smith',
  'https://lh3.googleusercontent.com/a/default-user=s96-c'),
  ('345678901234567890123', 'mike.johnson@example.com', 'Mike Johnson',
  'https://lh3.googleusercontent.com/a/default-user=s96-c');