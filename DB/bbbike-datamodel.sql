CREATE EXTENSION postgis;
CREATE EXTENSION postgis_topology;
CREATE EXTENSION fuzzystrmatch;
CREATE EXTENSION postgis_tiger_geocoder;

DROP TABLE city;
DROP TABLE cycleway;
DROP TABLE greenway;
DROP TABLE quality;
DROP TABLE trafficlight;
DROP TABLE node;
DROP TABLE way;

CREATE TABLE public.way(
	wayid bigserial,
	name name,
	type varchar,
	geometry geometry,
	attributes json,
	city bigint,
	nodes bigint[],
	CONSTRAINT wayid PRIMARY KEY (wayid)
);

CREATE TABLE public.city(
	name name,
	geometry geometry(MULTIPOLYGON),
	cityid bigserial,
	CONSTRAINT cityid PRIMARY KEY (cityid)
);

CREATE TABLE public.cycleway(
	cycleid bigserial,
	type varchar,
	geometry geometry,
	CONSTRAINT cycleid PRIMARY KEY (cycleid)
);

CREATE TABLE public.greenway(
	greenwayid bigserial,
	geometry geometry,
	type varchar,
	CONSTRAINT greenwayid PRIMARY KEY (greenwayid)
);

CREATE TABLE public.quality(
	qualityid bigserial,
	geometry geometry,
	type varchar,
	CONSTRAINT qualityid PRIMARY KEY (qualityid)
);

CREATE TABLE public.trafficlight(
	trafficlightid bigserial,
	geometry geometry,
	CONSTRAINT trafficlightid PRIMARY KEY (trafficlightid)
);

CREATE TABLE public.node(
	nodeid bigserial,
	geometry geometry,
	ways bigint[],
	neighbors bigint[],
	CONSTRAINT nodeid PRIMARY KEY (nodeid)
);

CREATE INDEX way_nodes_idx ON way USING gin (nodes);
CREATE INDEX nodes_neighbors_idx ON node USING gin (neighbors);
CREATE INDEX node_ways_idx ON node USING gin (ways);
CREATE INDEX way_gix ON way USING GIST (geometry);
CREATE INDEX cycleway_gix ON cycleway USING GIST (geometry);
CREATE INDEX node_gix ON node USING GIST (geometry);
CREATE INDEX greenway_gix ON greenway USING GIST (geometry);
CREATE INDEX quality_gix ON quality USING GIST (geometry);
