--
-- PostgreSQL database dump
--

-- Dumped from database version 15.8 (Debian 15.8-1.pgdg120+1)
-- Dumped by pg_dump version 15.8 (Debian 15.8-1.pgdg120+1)

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

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: backup; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.backup (
    id integer NOT NULL,
    bdd_id integer NOT NULL,
    backup_data bytea,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    file_size bigint,
    file_name character varying(255),
    status character varying(50) DEFAULT 'success'::character varying
);


ALTER TABLE public.backup OWNER TO admin;

--
-- Name: backup_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.backup_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.backup_id_seq OWNER TO admin;

--
-- Name: backup_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.backup_id_seq OWNED BY public.backup.id;


--
-- Name: databases; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.databases (
    id integer NOT NULL,
    dbname character varying(255) NOT NULL,
    dbport character varying(10) NOT NULL,
    username character varying(255) NOT NULL,
    password character varying(255) NOT NULL,
    type character varying(50) NOT NULL
);


ALTER TABLE public.databases OWNER TO admin;

--
-- Name: databases_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.databases_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.databases_id_seq OWNER TO admin;

--
-- Name: databases_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.databases_id_seq OWNED BY public.databases.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.users (
    id integer NOT NULL,
    login character varying(255) NOT NULL,
    password character varying(255) NOT NULL
);


ALTER TABLE public.users OWNER TO admin;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO admin;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: backup id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.backup ALTER COLUMN id SET DEFAULT nextval('public.backup_id_seq'::regclass);


--
-- Name: databases id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.databases ALTER COLUMN id SET DEFAULT nextval('public.databases_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: backup; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.backup (id, bdd_id, backup_data, created_at, file_size, file_name, status) FROM stdin;
\.


--
-- Data for Name: databases; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.databases (id, dbname, dbport, username, password, type) FROM stdin;
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.users (id, login, password) FROM stdin;
\.


--
-- Name: backup_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.backup_id_seq', 1, false);


--
-- Name: databases_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.databases_id_seq', 1, false);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.users_id_seq', 1, false);


--
-- Name: backup backup_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.backup
    ADD CONSTRAINT backup_pkey PRIMARY KEY (id);


--
-- Name: databases databases_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.databases
    ADD CONSTRAINT databases_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: backup backup_bdd_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.backup
    ADD CONSTRAINT backup_bdd_id_fkey FOREIGN KEY (bdd_id) REFERENCES public.databases(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

