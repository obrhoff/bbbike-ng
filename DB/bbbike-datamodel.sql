CREATE EXTENSION postgis;
CREATE EXTENSION postgis_topology;
CREATE EXTENSION fuzzystrmatch;
CREATE EXTENSION postgis_tiger_geocoder;
CREATE EXTENSION hstore;

DROP TABLE IF EXISTS city;

DROP TABLE IF EXISTS place;

DROP TABLE IF EXISTS path;
DROP TABLE IF EXISTS cyclepath;
DROP TABLE IF EXISTS greenpath;

DROP TABLE IF EXISTS quality;
DROP TABLE IF EXISTS trafficlight;
DROP TABLE IF EXISTS network;
DROP TABLE IF EXISTS node;

SELECT topology.DropTopology('path_topo');
SELECT topology.DropTopology('place_topo');

CREATE TABLE public.path(
	id bigserial,
	name name,
	type varchar,
	geometry geometry(linestring, 4326),
	city bigint,
	nodes bigint[],
	CONSTRAINT pathid PRIMARY KEY (id)
);

CREATE INDEX streetway_gix ON path USING GIST (geometry);

CREATE TABLE public.place(
	id bigserial,
	name name,
	type varchar,
	geometry geometry(point, 4326),
	city bigint,
	nodes bigint[],
	CONSTRAINT placeid PRIMARY KEY (id)
);

CREATE INDEX place_gix ON place USING GIST (geometry);

CREATE TABLE public.city(
	id bigserial,
	name name,
	geometry geometry(MULTIPOLYGON, 4326),
	CONSTRAINT cityid PRIMARY KEY (id)
);

CREATE TABLE public.cyclepath(
	id bigserial,
	type varchar,
	geometry geometry(linestring, 4326),
	CONSTRAINT cycleid PRIMARY KEY (id)
);
CREATE INDEX cyclepath_gix ON cyclepath USING GIST (geometry);

CREATE TABLE public.greenpath(
	id bigserial,
	geometry geometry(linestring, 4326),
	type varchar,
	CONSTRAINT greenwayid PRIMARY KEY (id)
);
CREATE INDEX greenpath_gix ON greenpath USING GIST (geometry);

CREATE TABLE public.quality(
	id bigserial,
	geometry geometry(linestring, 4326),
	type varchar,
	CONSTRAINT qualityid PRIMARY KEY (id)
);

CREATE INDEX quality_gix ON quality USING GIST (id);
CREATE TABLE public.trafficlight(
	id bigserial,
	geometry geometry(point, 4326),
	CONSTRAINT trafficlightid PRIMARY KEY (id)
);

CREATE TABLE public.node(
	id bigserial,
	geometry geometry(point, 4326),
	networks bigint[],
	neighbors bigint[],
	walkable bool,
	CONSTRAINT nodeid PRIMARY KEY (id)
);

CREATE INDEX node_gix ON node USING GIST (geometry);
CREATE INDEX nodes_neighbors_idx ON node USING gin (neighbors);
CREATE INDEX node_ways_idx ON node USING gin (networks);

CREATE TABLE public.network(
    id bigserial,
    type varchar,
    geometry geometry(linestring, 4326),
    nodes bigint[],
    attributes hstore NOT NULL DEFAULT,
    CONSTRAINT networkid PRIMARY KEY (networkid)
);

CREATE INDEX attributes_idx ON network USING gin(attributes);
CREATE INDEX network_idx ON network USING GIST (geometry);
CREATE INDEX network_nodes_idx ON network USING gin (nodes);

