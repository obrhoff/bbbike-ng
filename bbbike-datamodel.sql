-- Database generated with pgModeler (PostgreSQL Database Modeler).
-- PostgreSQL version: 9.1
-- Project Site: pgmodeler.com.br
-- Model Author: ---

SET check_function_bodies = false;
-- ddl-end --


-- Database creation must be done outside an multicommand file.
-- These commands were put in this file only for convenience.
-- -- object: new_database | type: DATABASE --
-- CREATE DATABASE new_database
-- ;
-- -- ddl-end --
-- 

-- object: public.streets | type: TABLE --
CREATE TABLE public.streets(
	streetid bigserial,
	name varchar,
	type varchar,
	streetpath path,
	CONSTRAINT streetid PRIMARY KEY (streetid)

);
-- ddl-end --
-- object: public.city | type: TABLE --
CREATE TABLE public.city(
	name name,
	bounds geometry(MULTIPOLYGON),
	cityid bigserial,
	CONSTRAINT cityid PRIMARY KEY (cityid)

);
-- ddl-end --

