CREATE TABLE public.answer (
    id integer NOT NULL,
    name character varying(100)
);


CREATE TABLE public.question (
    id integer NOT NULL,
    description character varying
);


CREATE TABLE public.survey (
    id integer NOT NULL,
    name text NOT NULL,
    status boolean DEFAULT true NOT NULL,
    rka text,
    rc_name character varying(255),
    adress character varying(255),
    question_id integer[] NOT NULL,
    answer_id integer[] NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    user_id integer NOT NULL
);