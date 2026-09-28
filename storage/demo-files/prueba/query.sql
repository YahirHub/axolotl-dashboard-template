CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    active BOOLEAN DEFAULT TRUE
);

SELECT id, name, email
FROM users
WHERE active = TRUE
ORDER BY name ASC;
