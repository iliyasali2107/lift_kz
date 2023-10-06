--
-- PostgreSQL database dump
--

-- Dumped from database version 15.1
-- Dumped by pg_dump version 15.1

-- Started on 2023-09-18 12:46:42 +06

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- TOC entry 224 (class 1255 OID 19570)
-- Name: f_isvalidemail(text); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.f_isvalidemail(text) RETURNS boolean
    LANGUAGE plpgsql
    AS $_$
BEGIN
    RETURN $1 ~ '^[a-zA-Z0-9.!#$%&''*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$';
END;
$_$;


ALTER FUNCTION public.f_isvalidemail(text) OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 220 (class 1259 OID 19553)
-- Name: answer; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.answer (
    id integer NOT NULL,
    name character varying(100)
);


ALTER TABLE public.answer OWNER TO postgres;

--
-- TOC entry 221 (class 1259 OID 19556)
-- Name: answer_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.answer ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.answer_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- TOC entry 218 (class 1259 OID 19544)
-- Name: question; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.question (
    id integer NOT NULL,
    description character varying
);


ALTER TABLE public.question OWNER TO postgres;

--
-- TOC entry 219 (class 1259 OID 19547)
-- Name: question_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.question ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.question_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- TOC entry 216 (class 1259 OID 19530)
-- Name: survey; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.survey (
    id integer NOT NULL,
    name text NOT NULL,
    status boolean DEFAULT true NOT NULL,
    rka text,
    rc_name character varying(255),
    adress character varying(255),
    question_id integer[] NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    user_id integer NOT NULL
);


ALTER TABLE public.survey OWNER TO postgres;

--
-- TOC entry 217 (class 1259 OID 19533)
-- Name: survey_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.survey ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.survey_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- TOC entry 222 (class 1259 OID 19574)
-- Name: user_question; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_question (
    id integer NOT NULL,
    user_id integer NOT NULL,
    question_id integer NOT NULL,
    answer_id integer
);


ALTER TABLE public.user_question OWNER TO postgres;

--
-- TOC entry 223 (class 1259 OID 19577)
-- Name: user_question_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.user_question ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.user_question_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- TOC entry 214 (class 1259 OID 19521)
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id integer NOT NULL,
    iin text NOT NULL,
    email character varying(255),
    bin character varying(255),
    is_manager boolean DEFAULT false NOT NULL,
    name character varying(255) NOT NULL
);


ALTER TABLE public.users OWNER TO postgres;

--
-- TOC entry 215 (class 1259 OID 19524)
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.users ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- TOC entry 3633 (class 0 OID 19553)
-- Dependencies: 220
-- Data for Name: answer; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.answer (id, name) FROM stdin;
1	да
2	нет
3	воздержусь
\.


--
-- TOC entry 3631 (class 0 OID 19544)
-- Dependencies: 218
-- Data for Name: question; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.question (id, description) FROM stdin;
\.


--
-- TOC entry 3629 (class 0 OID 19530)
-- Dependencies: 216
-- Data for Name: survey; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.survey (id, name, status, rka, rc_name, adress, question_id, created_at, user_id) FROM stdin;
\.


--
-- TOC entry 3635 (class 0 OID 19574)
-- Dependencies: 222
-- Data for Name: user_question; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.user_question (id, user_id, question_id, answer_id) FROM stdin;
\.


--
-- TOC entry 3627 (class 0 OID 19521)
-- Dependencies: 214
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, iin, email, bin, is_manager, name) FROM stdin;
\.


--
-- TOC entry 3642 (class 0 OID 0)
-- Dependencies: 221
-- Name: answer_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.answer_id_seq', 3, true);


--
-- TOC entry 3643 (class 0 OID 0)
-- Dependencies: 219
-- Name: question_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.question_id_seq', 1, false);


--
-- TOC entry 3644 (class 0 OID 0)
-- Dependencies: 217
-- Name: survey_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.survey_id_seq', 1, false);


--
-- TOC entry 3645 (class 0 OID 0)
-- Dependencies: 223
-- Name: user_question_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.user_question_id_seq', 1, false);


--
-- TOC entry 3646 (class 0 OID 0)
-- Dependencies: 215
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_id_seq', 5, true);


-- Completed on 2023-09-18 12:46:42 +06

--
-- PostgreSQL database dump complete
--

