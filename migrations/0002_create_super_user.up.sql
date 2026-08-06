INSERT INTO users (
    user_id, 
    name,
    last_name,
    password,
    email,
    role,
    verified,
    created_at,
    updated_at
    ) VALUES (
    gen_random_uuid(),
    'super',
    'admin',
    crypt('@Olecram_99', gen_salt('bf', 10)), 
    'superadmin@gmail.com',
    'super-admin',
    'yes',
    NOW(),
    NOW()
);