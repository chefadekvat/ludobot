INSERT INTO users (id, balance)
    VALUES ($1, $2)
ON CONFLICT (id) DO NOTHING;