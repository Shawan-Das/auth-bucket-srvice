-- SCHEMA: hrm
CREATE SCHEMA common;

CREATE TABLE common.users (
	user_id serial4 NOT NULL,
	user_name text NOT NULL,
	email text NOT NULL,
	phone text NOT NULL,
	pass text NOT NULL,
	pss_valid bool DEFAULT true NOT NULL,
	otp text NULL,
	otp_valid bool DEFAULT false NOT NULL,
	otp_exp timestamp NULL,
	"role" text NOT NULL
);