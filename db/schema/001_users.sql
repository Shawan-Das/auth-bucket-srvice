CREATE TABLE IF NOT EXISTS common.users (
    user_id serial4 PRIMARY KEY NOT NULL,
    user_name text NOT NULL,
    email text NOT NULL UNIQUE,
    phone text NOT NULL,
    pass text NOT NULL,
    pss_valid bool DEFAULT true NOT NULL,
    otp text NOT NULL DEFAULT '',
    otp_valid bool DEFAULT false NOT NULL,
    otp_exp timestamp NULL,
    role text NOT NULL,
    refresh_token text DEFAULT '',
    refresh_token_exp timestamp NULL,
    created_at timestamp DEFAULT now() NOT NULL,
    updated_at timestamp DEFAULT now() NOT NULL
);
