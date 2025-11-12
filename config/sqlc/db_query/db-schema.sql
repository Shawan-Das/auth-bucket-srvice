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

CREATE TABLE common.satcom_data (
	id serial4 NOT NULL,
	company text NOT NULL,
	category text NOT NULL,
	"type" text NOT NULL,
	"date" text NOT NULL,
	"time" text NOT NULL,
	db_port text NOT NULL,
	ui_port text NOT NULL,
	url text NOT NULL,
	ip text NOT NULL,
	status bool NOT NULL
);