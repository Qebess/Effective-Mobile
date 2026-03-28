
SET DateStyle = 'ISO, DMY';
CREATE TABLE IF NOT EXISTS subscription 
(
	id serial PRIMARY KEY NOT NULL,
  service_name varchar NOT NULL,
	price bigint NOT NULL,
	user_id uuid NOT NULL,
	start_date DATE
);

