CREATE EXTENSION postgis 
CREATE EXTENSION postgis_topology;
CREATE EXTENSION fuzzystrmatch;
CREATE EXTENSION postgis_tiger_geocoder;

CREATE TABLE public.streetpath(
	pathid bigserial,
	name varchar,
	type varchar,
	geometry path,
	attributes json,
	CONSTRAINT streetpathid PRIMARY KEY (pathid)

);

CREATE TABLE public.city(
	name name,
	bounds geometry(MULTIPOLYGON),
	cityid bigserial,
	CONSTRAINT cityid PRIMARY KEY (cityid)

);

CREATE TABLE public.cyclepath(
	pathid bigserial,
	type varchar,
	geometry path,
	attributes json,
	CONSTRAINT cyclepathid PRIMARY KEY (pathid)

);
