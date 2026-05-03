--
-- PostgreSQL database dump
--

\restrict H1yYwGq3cYyKiut79qglA98eWzRpghTifMObGazldYfFFaG41ptyAx6WvlmjTW1

-- Dumped from database version 18.2
-- Dumped by pg_dump version 18.2

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
-- Name: public; Type: SCHEMA; Schema: -; Owner: postgres
--

-- *not* creating schema, since initdb creates it


ALTER SCHEMA public OWNER TO postgres;

--
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: postgres
--

COMMENT ON SCHEMA public IS '';


--
-- Name: btree_gist; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS btree_gist WITH SCHEMA public;


--
-- Name: EXTENSION btree_gist; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION btree_gist IS 'support for indexing common datatypes in GiST';


--
-- Name: unaccent; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS unaccent WITH SCHEMA public;


--
-- Name: EXTENSION unaccent; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION unaccent IS 'text search dictionary that removes accents';


--
-- Name: status_type; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.status_type AS ENUM (
    'confirmed',
    'in_progress',
    'completed',
    'cancelled',
    'no_show'
);


ALTER TYPE public.status_type OWNER TO postgres;

--
-- Name: technician_level; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.technician_level AS ENUM (
    'fresher',
    'junior',
    'senior'
);


ALTER TYPE public.technician_level OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: appointments; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.appointments (
    id integer NOT NULL,
    dealership_id integer NOT NULL,
    service_id integer NOT NULL,
    bay_id integer NOT NULL,
    technician_id integer NOT NULL,
    customer_name character varying(100) NOT NULL,
    status public.status_type DEFAULT 'confirmed'::public.status_type NOT NULL,
    duration tstzrange NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT ck_appointments_duration_valid CHECK (((NOT isempty(duration)) AND (lower(duration) IS NOT NULL) AND (upper(duration) IS NOT NULL) AND (lower(duration) < upper(duration)) AND lower_inc(duration) AND (NOT upper_inc(duration))))
);


ALTER TABLE public.appointments OWNER TO postgres;

--
-- Name: appointments_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.appointments_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.appointments_id_seq OWNER TO postgres;

--
-- Name: appointments_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.appointments_id_seq OWNED BY public.appointments.id;


--
-- Name: dealerships; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.dealerships (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    open_time time without time zone NOT NULL,
    close_time time without time zone NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.dealerships OWNER TO postgres;

--
-- Name: dealerships_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.dealerships_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.dealerships_id_seq OWNER TO postgres;

--
-- Name: dealerships_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.dealerships_id_seq OWNED BY public.dealerships.id;


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO postgres;

--
-- Name: service_bay_types; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.service_bay_types (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.service_bay_types OWNER TO postgres;

--
-- Name: service_bay_types_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.service_bay_types_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.service_bay_types_id_seq OWNER TO postgres;

--
-- Name: service_bay_types_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.service_bay_types_id_seq OWNED BY public.service_bay_types.id;


--
-- Name: service_bays; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.service_bays (
    id integer NOT NULL,
    dealership_id integer NOT NULL,
    bay_type_id integer NOT NULL,
    name character varying(100) NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.service_bays OWNER TO postgres;

--
-- Name: service_bays_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.service_bays_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.service_bays_id_seq OWNER TO postgres;

--
-- Name: service_bays_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.service_bays_id_seq OWNED BY public.service_bays.id;


--
-- Name: service_requirements; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.service_requirements (
    service_id integer NOT NULL,
    skill_id integer NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.service_requirements OWNER TO postgres;

--
-- Name: services; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.services (
    id integer NOT NULL,
    required_bay_type_id integer NOT NULL,
    name character varying(100) NOT NULL,
    anticipated_minutes integer NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT services_anticipated_minutes_check CHECK ((anticipated_minutes > 0))
);


ALTER TABLE public.services OWNER TO postgres;

--
-- Name: services_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.services_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.services_id_seq OWNER TO postgres;

--
-- Name: services_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.services_id_seq OWNED BY public.services.id;


--
-- Name: skills; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.skills (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.skills OWNER TO postgres;

--
-- Name: skills_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.skills_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.skills_id_seq OWNER TO postgres;

--
-- Name: skills_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.skills_id_seq OWNED BY public.skills.id;


--
-- Name: technician_skills; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.technician_skills (
    technician_id integer NOT NULL,
    skill_id integer NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.technician_skills OWNER TO postgres;

--
-- Name: technicians; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.technicians (
    id integer NOT NULL,
    dealership_id integer NOT NULL,
    name character varying(50) NOT NULL,
    level public.technician_level NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    inactive_since timestamp with time zone,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.technicians OWNER TO postgres;

--
-- Name: technicians_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.technicians_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.technicians_id_seq OWNER TO postgres;

--
-- Name: technicians_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.technicians_id_seq OWNED BY public.technicians.id;


--
-- Name: appointments id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.appointments ALTER COLUMN id SET DEFAULT nextval('public.appointments_id_seq'::regclass);


--
-- Name: dealerships id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.dealerships ALTER COLUMN id SET DEFAULT nextval('public.dealerships_id_seq'::regclass);


--
-- Name: service_bay_types id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_bay_types ALTER COLUMN id SET DEFAULT nextval('public.service_bay_types_id_seq'::regclass);


--
-- Name: service_bays id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_bays ALTER COLUMN id SET DEFAULT nextval('public.service_bays_id_seq'::regclass);


--
-- Name: services id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.services ALTER COLUMN id SET DEFAULT nextval('public.services_id_seq'::regclass);


--
-- Name: skills id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.skills ALTER COLUMN id SET DEFAULT nextval('public.skills_id_seq'::regclass);


--
-- Name: technicians id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.technicians ALTER COLUMN id SET DEFAULT nextval('public.technicians_id_seq'::regclass);


--
-- Data for Name: appointments; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.appointments (id, dealership_id, service_id, bay_id, technician_id, customer_name, status, duration, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: dealerships; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.dealerships (id, name, open_time, close_time, created_at, updated_at) FROM stdin;
2	Dealership B	08:00:00	18:00:00	2026-05-01 18:08:29.392602+00	2026-05-01 18:08:29.392602+00
3	Dealership C	08:00:00	18:00:00	2026-05-01 18:08:34.132716+00	2026-05-01 18:08:34.132716+00
1	Dealership A	08:00:00	18:00:00	2026-05-02 04:00:04.616489+00	2026-05-02 04:00:04.616489+00
\.


--
-- Data for Name: schema_migrations; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.schema_migrations (version, dirty) FROM stdin;
10	f
\.


--
-- Data for Name: service_bay_types; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.service_bay_types (id, name, created_at, updated_at) FROM stdin;
1	Standard Bay	2026-05-01 18:11:03.756118+00	2026-05-01 18:11:03.756118+00
3	EV Bay	2026-05-01 18:11:22.211492+00	2026-05-01 18:11:22.211492+00
2	Lift Bay	2026-05-01 18:11:13.450893+00	2026-05-01 18:15:04.443999+00
5	Wash Bay	2026-05-01 18:17:03.314641+00	2026-05-01 18:17:03.314641+00
\.


--
-- Data for Name: service_bays; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.service_bays (id, dealership_id, bay_type_id, name, is_active, created_at, updated_at) FROM stdin;
5	2	1	SB 1	t	2026-05-02 02:09:25.6178+00	2026-05-02 02:09:25.6178+00
6	2	2	SB 2	t	2026-05-02 02:09:31.597603+00	2026-05-02 02:09:31.597603+00
7	2	3	SB 3	t	2026-05-02 02:09:39.309611+00	2026-05-02 02:09:39.309611+00
8	3	1	SB 1	t	2026-05-02 02:10:02.639533+00	2026-05-02 02:10:02.639533+00
9	3	2	SB 2	t	2026-05-02 02:10:08.910804+00	2026-05-02 02:10:08.910804+00
10	3	3	SB 3	t	2026-05-02 02:10:15.530119+00	2026-05-02 02:10:15.530119+00
12	1	1	SB 1	t	2026-05-02 04:20:50.477086+00	2026-05-02 04:20:50.477086+00
13	1	2	SB 2	t	2026-05-02 04:20:57.861504+00	2026-05-02 04:20:57.861504+00
14	1	3	SB 3	t	2026-05-02 04:21:05.086761+00	2026-05-02 04:21:05.086761+00
15	1	2	SB 4	t	2026-05-02 12:15:02.540319+00	2026-05-02 12:15:02.540319+00
\.


--
-- Data for Name: service_requirements; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.service_requirements (service_id, skill_id, created_at) FROM stdin;
1	12	2026-05-02 06:59:39.381945+00
2	4	2026-05-02 07:01:27.578531+00
2	5	2026-05-02 07:01:27.578531+00
3	9	2026-05-02 07:04:42.419447+00
3	10	2026-05-02 07:04:42.419447+00
2	6	2026-05-02 07:09:31.812236+00
5	7	2026-05-02 07:11:18.869974+00
5	8	2026-05-02 07:11:18.869974+00
6	3	2026-05-02 07:14:32.16509+00
6	2	2026-05-02 07:14:32.16509+00
\.


--
-- Data for Name: services; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.services (id, required_bay_type_id, name, anticipated_minutes, created_at, updated_at) FROM stdin;
1	1	Routine Maintenance	90	2026-05-02 06:55:27.887893+00	2026-05-02 06:55:27.887893+00
2	2	Brake Service	120	2026-05-02 07:00:42.204049+00	2026-05-02 07:00:42.204049+00
3	3	Hybrid/EV Service	240	2026-05-02 07:03:29.71826+00	2026-05-02 07:03:29.71826+00
5	1	Transmission Service	240	2026-05-02 07:05:35.907792+00	2026-05-02 07:10:25.258618+00
6	2	Engine Repair	240	2026-05-02 07:13:50.918089+00	2026-05-02 07:13:50.918089+00
\.


--
-- Data for Name: skills; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.skills (id, name, created_at, updated_at) FROM stdin;
2	Engine Diagnostics	2026-05-02 01:55:03.463517+00	2026-05-02 01:55:03.463517+00
3	Mechanical Overhaul	2026-05-02 01:55:43.037686+00	2026-05-02 01:55:43.037686+00
4	Brake System Maintenance	2026-05-02 01:56:08.630094+00	2026-05-02 01:56:08.630094+00
5	ABS Troubleshooting	2026-05-02 01:56:20.64676+00	2026-05-02 01:56:20.64676+00
6	Hydraulic Systems	2026-05-02 01:56:32.667516+00	2026-05-02 01:56:32.667516+00
7	Gearbox Overhaul	2026-05-02 01:57:07.168683+00	2026-05-02 01:57:07.168683+00
8	Transmission Diagnostics	2026-05-02 01:57:23.996088+00	2026-05-02 01:57:23.996088+00
9	Electric Motor Diagnostics	2026-05-02 01:58:10.556308+00	2026-05-02 01:58:10.556308+00
10	Battery Cell Management	2026-05-02 01:58:33.086296+00	2026-05-02 01:58:33.086296+00
12	Fluid & Filter Replacement	2026-05-02 04:19:59.675+00	2026-05-02 04:19:59.675+00
\.


--
-- Data for Name: technician_skills; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.technician_skills (technician_id, skill_id, created_at) FROM stdin;
1	5	2026-05-02 04:22:33.98175+00
1	4	2026-05-02 04:22:33.98175+00
1	12	2026-05-02 04:22:33.98175+00
1	6	2026-05-02 04:22:33.98175+00
2	10	2026-05-02 04:31:53.331794+00
2	9	2026-05-02 04:31:53.331794+00
2	2	2026-05-02 04:31:53.331794+00
2	3	2026-05-02 04:31:53.331794+00
4	10	2026-05-02 04:32:15.996874+00
4	9	2026-05-02 04:32:15.996874+00
4	2	2026-05-02 04:32:15.996874+00
4	7	2026-05-02 04:32:15.996874+00
4	3	2026-05-02 04:32:15.996874+00
4	8	2026-05-02 04:32:15.996874+00
3	5	2026-05-02 04:32:30.966345+00
3	4	2026-05-02 04:32:30.966345+00
3	12	2026-05-02 04:32:30.966345+00
3	6	2026-05-02 04:32:30.966345+00
6	5	2026-05-02 04:33:12.80222+00
6	4	2026-05-02 04:33:12.80222+00
6	12	2026-05-02 04:33:12.80222+00
6	6	2026-05-02 04:33:12.80222+00
5	10	2026-05-02 04:33:22.409224+00
5	9	2026-05-02 04:33:22.409224+00
5	2	2026-05-02 04:33:22.409224+00
5	7	2026-05-02 04:33:22.409224+00
5	3	2026-05-02 04:33:22.409224+00
5	8	2026-05-02 04:33:22.409224+00
1	7	2026-05-02 04:36:05.5246+00
1	8	2026-05-02 04:36:05.5246+00
8	4	2026-05-02 09:17:16.24747+00
8	5	2026-05-02 09:17:16.24747+00
8	6	2026-05-02 09:17:16.24747+00
\.


--
-- Data for Name: technicians; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.technicians (id, dealership_id, name, level, is_active, inactive_since, created_at, updated_at) FROM stdin;
2	1	John	junior	t	\N	2026-05-02 04:04:33.161238+00	2026-05-02 04:04:33.161238+00
3	2	Jack	junior	t	\N	2026-05-02 04:04:43.277568+00	2026-05-02 04:04:43.277568+00
4	2	Brown	senior	t	\N	2026-05-02 04:04:54.682306+00	2026-05-02 04:04:54.682306+00
5	3	Jason	fresher	t	\N	2026-05-02 04:05:40.646093+00	2026-05-02 04:05:40.646093+00
6	3	Henry	junior	t	\N	2026-05-02 04:06:08.62667+00	2026-05-02 04:06:08.62667+00
1	1	Danny	junior	t	\N	2026-05-02 04:02:19.582344+00	2026-05-02 04:10:57.03725+00
8	1	Victor	junior	t	\N	2026-05-02 09:16:07.589352+00	2026-05-02 09:16:07.589352+00
\.


--
-- Name: appointments_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.appointments_id_seq', 76, true);


--
-- Name: dealerships_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.dealerships_id_seq', 5, true);


--
-- Name: service_bay_types_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.service_bay_types_id_seq', 6, true);


--
-- Name: service_bays_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.service_bays_id_seq', 15, true);


--
-- Name: services_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.services_id_seq', 6, true);


--
-- Name: skills_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.skills_id_seq', 12, true);


--
-- Name: technicians_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.technicians_id_seq', 8, true);


--
-- Name: appointments appointments_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.appointments
    ADD CONSTRAINT appointments_pkey PRIMARY KEY (id);


--
-- Name: dealerships dealerships_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.dealerships
    ADD CONSTRAINT dealerships_pkey PRIMARY KEY (id);


--
-- Name: appointments ex_appointments_no_overlap_bay; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.appointments
    ADD CONSTRAINT ex_appointments_no_overlap_bay EXCLUDE USING gist (bay_id WITH =, duration WITH &&) WHERE (((bay_id IS NOT NULL) AND (status <> 'cancelled'::public.status_type)));


--
-- Name: appointments ex_appointments_no_overlap_technician; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.appointments
    ADD CONSTRAINT ex_appointments_no_overlap_technician EXCLUDE USING gist (technician_id WITH =, duration WITH &&) WHERE (((technician_id IS NOT NULL) AND (status <> 'cancelled'::public.status_type)));


--
-- Name: service_requirements pk_service_requirements; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_requirements
    ADD CONSTRAINT pk_service_requirements PRIMARY KEY (service_id, skill_id);


--
-- Name: technician_skills pk_technician_skills; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.technician_skills
    ADD CONSTRAINT pk_technician_skills PRIMARY KEY (technician_id, skill_id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: service_bay_types service_bay_types_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_bay_types
    ADD CONSTRAINT service_bay_types_pkey PRIMARY KEY (id);


--
-- Name: service_bays service_bays_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_bays
    ADD CONSTRAINT service_bays_pkey PRIMARY KEY (id);


--
-- Name: services services_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.services
    ADD CONSTRAINT services_pkey PRIMARY KEY (id);


--
-- Name: skills skills_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.skills
    ADD CONSTRAINT skills_pkey PRIMARY KEY (id);


--
-- Name: technicians technicians_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.technicians
    ADD CONSTRAINT technicians_pkey PRIMARY KEY (id);


--
-- Name: service_bays ux_service_bays_dealership_name; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_bays
    ADD CONSTRAINT ux_service_bays_dealership_name UNIQUE (dealership_id, name);


--
-- Name: ix_appointments_bay_duration_gist; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX ix_appointments_bay_duration_gist ON public.appointments USING gist (bay_id, duration) WHERE ((bay_id IS NOT NULL) AND (status <> 'cancelled'::public.status_type));


--
-- Name: ix_appointments_dealership_duration_gist; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX ix_appointments_dealership_duration_gist ON public.appointments USING gist (dealership_id, duration);


--
-- Name: ix_appointments_status_duration_gist; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX ix_appointments_status_duration_gist ON public.appointments USING gist (status, duration);


--
-- Name: ix_appointments_technician_duration_gist; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX ix_appointments_technician_duration_gist ON public.appointments USING gist (technician_id, duration) WHERE ((technician_id IS NOT NULL) AND (status <> 'cancelled'::public.status_type));


--
-- Name: ix_service_bays_active; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX ix_service_bays_active ON public.service_bays USING btree (dealership_id, is_active);


--
-- Name: ix_service_bays_dealership_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX ix_service_bays_dealership_id ON public.service_bays USING btree (dealership_id);


--
-- Name: ix_service_bays_type; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX ix_service_bays_type ON public.service_bays USING btree (bay_type_id);


--
-- Name: ix_service_requirements_service_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX ix_service_requirements_service_id ON public.service_requirements USING btree (service_id);


--
-- Name: ix_service_requirements_skill_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX ix_service_requirements_skill_id ON public.service_requirements USING btree (skill_id);


--
-- Name: ix_technician_skills_skill_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX ix_technician_skills_skill_id ON public.technician_skills USING btree (skill_id);


--
-- Name: ix_technician_skills_technician_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX ix_technician_skills_technician_id ON public.technician_skills USING btree (technician_id);


--
-- Name: ix_technicians_dealership_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX ix_technicians_dealership_id ON public.technicians USING btree (dealership_id);


--
-- Name: ux_dealerships_name; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX ux_dealerships_name ON public.dealerships USING btree (name);


--
-- Name: ux_service_bay_types_name; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX ux_service_bay_types_name ON public.service_bay_types USING btree (name);


--
-- Name: ux_services_name; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX ux_services_name ON public.services USING btree (name);


--
-- Name: ux_skills_name; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX ux_skills_name ON public.skills USING btree (name);


--
-- Name: appointments appointments_bay_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.appointments
    ADD CONSTRAINT appointments_bay_id_fkey FOREIGN KEY (bay_id) REFERENCES public.service_bays(id) ON DELETE SET NULL;


--
-- Name: appointments appointments_dealership_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.appointments
    ADD CONSTRAINT appointments_dealership_id_fkey FOREIGN KEY (dealership_id) REFERENCES public.dealerships(id) ON DELETE CASCADE;


--
-- Name: appointments appointments_service_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.appointments
    ADD CONSTRAINT appointments_service_id_fkey FOREIGN KEY (service_id) REFERENCES public.services(id) ON DELETE SET NULL;


--
-- Name: appointments appointments_technician_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.appointments
    ADD CONSTRAINT appointments_technician_id_fkey FOREIGN KEY (technician_id) REFERENCES public.technicians(id) ON DELETE SET NULL;


--
-- Name: service_bays service_bays_bay_type_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_bays
    ADD CONSTRAINT service_bays_bay_type_id_fkey FOREIGN KEY (bay_type_id) REFERENCES public.service_bay_types(id) ON DELETE SET NULL;


--
-- Name: service_bays service_bays_dealership_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_bays
    ADD CONSTRAINT service_bays_dealership_id_fkey FOREIGN KEY (dealership_id) REFERENCES public.dealerships(id) ON DELETE CASCADE;


--
-- Name: service_requirements service_requirements_service_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_requirements
    ADD CONSTRAINT service_requirements_service_id_fkey FOREIGN KEY (service_id) REFERENCES public.services(id) ON DELETE CASCADE;


--
-- Name: service_requirements service_requirements_skill_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_requirements
    ADD CONSTRAINT service_requirements_skill_id_fkey FOREIGN KEY (skill_id) REFERENCES public.skills(id) ON DELETE CASCADE;


--
-- Name: services services_required_bay_type_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.services
    ADD CONSTRAINT services_required_bay_type_id_fkey FOREIGN KEY (required_bay_type_id) REFERENCES public.service_bay_types(id) ON DELETE SET NULL;


--
-- Name: technician_skills technician_skills_skill_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.technician_skills
    ADD CONSTRAINT technician_skills_skill_id_fkey FOREIGN KEY (skill_id) REFERENCES public.skills(id) ON DELETE CASCADE;


--
-- Name: technician_skills technician_skills_technician_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.technician_skills
    ADD CONSTRAINT technician_skills_technician_id_fkey FOREIGN KEY (technician_id) REFERENCES public.technicians(id) ON DELETE CASCADE;


--
-- Name: technicians technicians_dealership_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.technicians
    ADD CONSTRAINT technicians_dealership_id_fkey FOREIGN KEY (dealership_id) REFERENCES public.dealerships(id) ON DELETE CASCADE;


--
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: postgres
--

REVOKE USAGE ON SCHEMA public FROM PUBLIC;


--
-- PostgreSQL database dump complete
--

\unrestrict H1yYwGq3cYyKiut79qglA98eWzRpghTifMObGazldYfFFaG41ptyAx6WvlmjTW1

