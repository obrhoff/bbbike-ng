CREATE EXTENSION postgis;
CREATE EXTENSION postgis_topology;
CREATE EXTENSION fuzzystrmatch;
CREATE EXTENSION postgis_tiger_geocoder;

CREATE TABLE public.streetpaths(
	streetpathid bigserial,
	name varchar,
	type varchar,
	path geometry,
	attributes json,
	city bigserial,
	CONSTRAINT streetpathid PRIMARY KEY (streetpathid)

);

CREATE TABLE public.city(
	name name,
	bounds geometry(MULTIPOLYGON),
	cityid bigserial,
	CONSTRAINT cityid PRIMARY KEY (cityid)
);

CREATE TABLE public.cyclepaths(
	cyclepathid bigserial,
	type varchar,
	path geometry,
	CONSTRAINT cyclepathid PRIMARY KEY (cyclepathid)

);

CREATE TABLE public.greenways(
	greenwayid bigserial,
	path geometry,
	type varchar,
	CONSTRAINT greenwayid PRIMARY KEY (greenwayid)

);

CREATE TABLE public.qualitys(
	qualityid bigserial,
	path geometry,
	type varchar,
	CONSTRAINT qualityid PRIMARY KEY (qualityid)

);

CREATE TABLE public.trafficlights(
	trafficlightid bigserial,
	path geometry,
	CONSTRAINT trafficlightid PRIMARY KEY (trafficlightid)

);
