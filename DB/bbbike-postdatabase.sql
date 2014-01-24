SET search_path = topology,public;

SELECT topology.CreateTopology('way_topo', 4326);
SELECT topology.AddTopoGeometryColumn('way_topo', 'public', 'way', 'topo_geom', 'LINESTRING');
UPDATE way SET topo_geom = topology.toTopoGeom(geometry, 'way_topo', 1, 0.0);

CREATE TABLE public.way_cut(
    wayid bigserial,
    name name,
    type varchar,
    geometry geometry(linestring, 4326),
    attributes json,
    city bigint,
    nodes bigint[],
    CONSTRAINT waycutid PRIMARY KEY (waycutid)
);

CREATE INDEX way_cut_gix ON way_cut USING GIST (geometry);

insert into way_cut(name, type, geometry) SELECT w.name, w.type, e.geom
FROM way_topo.edge e,
     way_topo.relation rel,
     way w
WHERE e.edge_id = rel.element_id
AND rel.topogeo_id = (w.topo_geom).id;

INSERT INTO node (geometry)
SELECT (ST_DumpPoints(ST_Intersection(s1.geometry, s2.geometry))).geom as geometry
FROM way_cut s1, way_cut s2
WHERE ST_Intersects(s2.geometry, s1.geometry) AND s1.wayid != s2.wayid;

DELETE FROM node WHERE nodeid IN
( SELECT mt1.nodeid FROM node mt1, node mt2 WHERE mt1.geometry ~=  mt2.geometry AND mt1.nodeid < mt2.nodeid );

UPDATE node
SET ways = subquery.ways
FROM ( SELECT node.nodeid, array_agg(way_cut.wayid) AS ways FROM way_cut, node WHERE ST_Intersects(way_cut.geometry, node.geometry) GROUP BY node.nodeid ) AS subquery
WHERE subquery.nodeid = node.nodeid;

UPDATE way_cut
SET nodes = subquery.nodes
FROM ( SELECT way_cut.wayid, array_agg(node.nodeid) AS nodes FROM node, way_cut WHERE ST_Intersects(way_cut.geometry, node.geometry) GROUP BY way_cut.wayid ) AS subquery
WHERE subquery.wayid = way_cut.wayid;

UPDATE node
SET neighbors = subquery.neighbors
FROM ( SELECT node1.nodeid, array_agg(node2.nodeid) AS neighbors FROM node node1, node node2 WHERE node1.ways && node2.ways AND node1.nodeid != node2.nodeid GROUP BY node1.nodeid ) AS subquery
WHERE subquery.nodeid = node.nodeid;

DROP TABLE way;
alter table way_cut rename to way;

UPDATE node SET walkable = true WHERE array_length(neighbors, 1) > 1;
UPDATE node SET walkable = false WHERE array_length(neighbors, 1) < 2;