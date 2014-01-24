CREATE EXTENSION postgis;
CREATE EXTENSION postgis_topology;
CREATE EXTENSION fuzzystrmatch;
CREATE EXTENSION postgis_tiger_geocoder;

DROP TABLE IF EXISTS city;
DROP TABLE IF EXISTS cycleway;
DROP TABLE IF EXISTS greenway;
DROP TABLE IF EXISTS quality;
DROP TABLE IF EXISTS trafficlight;
DROP TABLE IF EXISTS node;
DROP TABLE IF EXISTS way;
DROP TABLE IF EXISTS place;

CREATE TABLE public.way(
	wayid bigserial,
	name name,
	type varchar,
	geometry geometry(linestring, 4326),
	attributes json,
	city bigint,
	nodes bigint[],
	CONSTRAINT wayid PRIMARY KEY (wayid)
);

CREATE TABLE public.place(
	placeid bigserial,
	name name,
	type varchar,
	geometry geometry(point, 4326),
	attributes json,
	city bigint,
	nodes bigint[],
	CONSTRAINT placeid PRIMARY KEY (placeid)
);

CREATE TABLE public.city(
	name name,
	geometry geometry(MULTIPOLYGON, 4326),
	cityid bigserial,
	CONSTRAINT cityid PRIMARY KEY (cityid)
);

CREATE TABLE public.cycleway(
	cycleid bigserial,
	type varchar,
	geometry geometry(linestring, 4326),
	CONSTRAINT cycleid PRIMARY KEY (cycleid)
);

CREATE TABLE public.greenway(
	greenwayid bigserial,
	geometry geometry(linestring, 4326),
	type varchar,
	CONSTRAINT greenwayid PRIMARY KEY (greenwayid)
);

CREATE TABLE public.quality(
	qualityid bigserial,
	geometry geometry(linestring, 4326),
	type varchar,
	CONSTRAINT qualityid PRIMARY KEY (qualityid)
);

CREATE TABLE public.trafficlight(
	trafficlightid bigserial,
	geometry geometry(linestring, 4326),
	CONSTRAINT trafficlightid PRIMARY KEY (trafficlightid)
);

CREATE TABLE public.node(
	nodeid bigserial,
	geometry geometry(point, 4326),
	ways bigint[],
	neighbors bigint[],
	CONSTRAINT nodeid PRIMARY KEY (nodeid)
);

CREATE INDEX place_gix ON place USING GIST (geometry);
CREATE INDEX way_nodes_idx ON way USING gin (nodes);
CREATE INDEX nodes_neighbors_idx ON node USING gin (neighbors);
CREATE INDEX node_ways_idx ON node USING gin (ways);
CREATE INDEX way_gix ON way USING GIST (geometry);
CREATE INDEX cycleway_gix ON cycleway USING GIST (geometry);
CREATE INDEX node_gix ON node USING GIST (geometry);
CREATE INDEX greenway_gix ON greenway USING GIST (geometry);
CREATE INDEX quality_gix ON quality USING GIST (geometry);
