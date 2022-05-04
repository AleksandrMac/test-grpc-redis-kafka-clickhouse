CREATE TABLE public.user
(
    id SERIAL,
	email VARCHAR(150) NOT NULL,
    CONSTRAINT user_pkey PRIMARY KEY (id),
    CONSTRAINT user_email_key UNIQUE (email)
);