SET search_path = topology,public;

SELECT topology.CreateTopology('way_topo', 4326);
SELECT topology.AddTopoGeometryColumn('way_topo', 'public', 'way', 'topo_geom', 'LINESTRING');
UPDATE way SET topo_geom = topology.toTopoGeom(geometry, 'way_topo', 1, 0.0);


insert into way_segment(wayid, geometry) SELECT w.wayid, e.geom
FROM way_topo.edge e,
     way_topo.relation rel,
     way w
WHERE e.edge_id = rel.element_id
AND rel.topogeo_id = (w.topo_geom).id;

INSERT INTO node (geometry)
SELECT (ST_DumpPoints(ST_Intersection(s1.geometry, s2.geometry))).geom as geometry
FROM way_segment s1, way_segment s2
WHERE ST_Intersects(s2.geometry, s1.geometry) AND s1.waysegmentid != s2.waysegmentid;

DELETE FROM node WHERE nodeid IN
( SELECT mt1.nodeid FROM node mt1, node mt2 WHERE mt1.geometry ~=  mt2.geometry AND mt1.nodeid < mt2.nodeid );

UPDATE node
SET waysegments = subquery.ways
FROM ( SELECT node.nodeid, array_agg(way_segment.waysegmentid) AS ways FROM way_segment, node WHERE ST_Intersects(way_segment.geometry, node.geometry) GROUP BY node.nodeid ) AS subquery
WHERE subquery.nodeid = node.nodeid;

UPDATE way_segment
SET nodes = subquery.nodes
FROM ( SELECT way_segment.waysegmentid, array_agg(node.nodeid) AS nodes FROM node, way_segment WHERE ST_Intersects(way_segment.geometry, node.geometry) GROUP BY way_segment.waysegmentid ) AS subquery
WHERE subquery.waysegmentid = way_segment.waysegmentid;

UPDATE node
SET neighbors = subquery.neighbors
FROM ( SELECT node1.nodeid, array_agg(node2.nodeid) AS neighbors FROM node node1, node node2 WHERE node1.waysegments && node2.waysegments AND node1.nodeid != node2.nodeid GROUP BY node1.nodeid ) AS subquery
WHERE subquery.nodeid = node.nodeid;

DROP TABLE way;
alter table way_cut rename to way;

UPDATE node SET walkable = true WHERE array_length(neighbors, 1) > 1;
UPDATE node SET walkable = false WHERE array_length(neighbors, 1) < 2;