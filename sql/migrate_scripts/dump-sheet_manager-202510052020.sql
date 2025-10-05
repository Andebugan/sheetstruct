--
-- PostgreSQL database cluster dump
--

-- Started on 2025-10-05 20:20:31

\restrict YysabLRB55WiI5JuaJlGf3cFazH4ITr0nOlLECB2mAqSlcMOqzBxf9SJSCACg7b

SET default_transaction_read_only = off;

SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;

--
-- Roles
--

CREATE ROLE postgres;
ALTER ROLE postgres WITH SUPERUSER INHERIT CREATEROLE CREATEDB LOGIN REPLICATION BYPASSRLS;

--
-- User Configurations
--








\unrestrict YysabLRB55WiI5JuaJlGf3cFazH4ITr0nOlLECB2mAqSlcMOqzBxf9SJSCACg7b

--
-- Databases
--

--
-- Database "template1" dump
--

\connect template1

--
-- PostgreSQL database dump
--

\restrict dzDNbjrVnh5qi0ZIHphCT21yopLJufZ3sXcvqDgUSfM5DKwwuz6HnYEYROT0itd

-- Dumped from database version 18.0
-- Dumped by pg_dump version 18.0

-- Started on 2025-10-05 20:20:31

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

-- Completed on 2025-10-05 20:20:31

--
-- PostgreSQL database dump complete
--

\unrestrict dzDNbjrVnh5qi0ZIHphCT21yopLJufZ3sXcvqDgUSfM5DKwwuz6HnYEYROT0itd

--
-- Database "sheet_manager" dump
--

--
-- PostgreSQL database dump
--

\restrict hbkfsqWbne7FxCf5WXnkx0eRV6v6cyKqTgyOFTevfuI8Oiwja6OoSf4SATvZEBp

-- Dumped from database version 18.0
-- Dumped by pg_dump version 18.0

-- Started on 2025-10-05 20:20:31

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- TOC entry 5003 (class 1262 OID 16388)
-- Name: sheet_manager; Type: DATABASE; Schema: -; Owner: postgres
--

CREATE DATABASE sheet_manager WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'Russian_Russia.1251';


ALTER DATABASE sheet_manager OWNER TO postgres;

\unrestrict hbkfsqWbne7FxCf5WXnkx0eRV6v6cyKqTgyOFTevfuI8Oiwja6OoSf4SATvZEBp
\connect sheet_manager
\restrict hbkfsqWbne7FxCf5WXnkx0eRV6v6cyKqTgyOFTevfuI8Oiwja6OoSf4SATvZEBp

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

-- Completed on 2025-10-05 20:20:32

--
-- PostgreSQL database dump complete
--

\unrestrict hbkfsqWbne7FxCf5WXnkx0eRV6v6cyKqTgyOFTevfuI8Oiwja6OoSf4SATvZEBp

-- Completed on 2025-10-05 20:20:32

--
-- PostgreSQL database cluster dump complete
--

