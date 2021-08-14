--
-- PostgreSQL database dump
--

-- Dumped from database version 13.3 (Debian 13.3-1.pgdg100+1)
-- Dumped by pg_dump version 13.3 (Debian 13.3-1.pgdg100+1)

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
-- Name: accounts; Type: TABLE; Schema: public; Owner: postgres
--
CREATE DATABASE ctf0;

\connect ctf0

CREATE TABLE public.accounts (
    username text NOT NULL,
    email text NOT NULL,
    password text NOT NULL,
    entiere bigint NOT NULL,
    decimale integer NOT NULL
);


ALTER TABLE public.accounts OWNER TO postgres;

--
-- Data for Name: accounts; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.accounts (username, email, password, entiere, decimale) FROM stdin;
admin	admin@admin.com	0000	9999998984	85
\.


--
-- Name: accounts accounts_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.accounts
    ADD CONSTRAINT "accounts_pkey" PRIMARY KEY (username);


--
-- PostgreSQL database dump complete
--

