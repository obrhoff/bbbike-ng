CREATE EXTENSION postgis;
CREATE EXTENSION postgis_topology;
CREATE EXTENSION fuzzystrmatch;
CREATE EXTENSION postgis_tiger_geocoder;
CREATE EXTENSION hstore;

DROP TABLE IF EXISTS city;

DROP TABLE IF EXISTS place;

DROP TABLE IF EXISTS streetway;
DROP TABLE IF EXISTS cycleway;
DROP TABLE IF EXISTS greenway;

DROP TABLE IF EXISTS quality;
DROP TABLE IF EXISTS trafficlight;
DROP TABLE IF EXISTS network;
DROP TABLE IF EXISTS node;

SELECT topology.DropTopology('network_topo');
SELECT topology.DropTopology('place_topo');

CREATE TABLE public.streetway(
	streetwayid bigserial,
	name name,
	type varchar,
	geometry geometry(linestring, 4326),
	city bigint,
	nodes bigint[],
	CONSTRAINT streetwayid PRIMARY KEY (streetwayid)
);

CREATE INDEX streetway_gix ON streetway USING GIST (geometry);

CREATE TABLE public.place(
	placeid bigserial,
	name name,
	type varchar,
	geometry geometry(point, 4326),
	city bigint,
	nodes bigint[],
	CONSTRAINT placeid PRIMARY KEY (placeid)
);

CREATE INDEX place_gix ON place USING GIST (geometry);

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
CREATE INDEX cycleway_gix ON cycleway USING GIST (geometry);

CREATE TABLE public.greenway(
	greenwayid bigserial,
	geometry geometry(linestring, 4326),
	type varchar,
	CONSTRAINT greenwayid PRIMARY KEY (greenwayid)
);
CREATE INDEX greenway_gix ON greenway USING GIST (geometry);

CREATE TABLE public.quality(
	qualityid bigserial,
	geometry geometry(linestring, 4326),
	type varchar,
	CONSTRAINT qualityid PRIMARY KEY (qualityid)
);

CREATE INDEX quality_gix ON quality USING GIST (geometry);
CREATE TABLE public.trafficlight(
	trafficlightid bigserial,
	geometry geometry(point, 4326),
	CONSTRAINT trafficlightid PRIMARY KEY (trafficlightid)
);

CREATE TABLE public.node(
	nodeid bigserial,
	geometry geometry(point, 4326),
	networks bigint[],
	neighbors bigint[],
	walkable bool,
	trafficlight bool,
	CONSTRAINT nodeid PRIMARY KEY (nodeid)
);

CREATE INDEX node_gix ON node USING GIST (geometry);
CREATE INDEX nodes_neighbors_idx ON node USING gin (neighbors);
CREATE INDEX node_ways_idx ON node USING gin (networks);

CREATE TABLE public.network(
    networkid bigserial,
    foreignid bigint,
    type varchar,
    geometry geometry(linestring, 4326),
    nodes bigint[],
    CONSTRAINT networkid PRIMARY KEY (networkid)
);

CREATE INDEX network_gix ON network USING GIST (geometry);
CREATE INDEX network_nodes_idx ON network USING gin (nodes);

