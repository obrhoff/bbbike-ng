INSERT INTO node (geometry)
SELECT (ST_DumpPoints(ST_Intersection(s1.geometry, s2.geometry))).geom as geometry FROM way s1, way s2 WHERE ST_Intersects(s2.geometry, s1.geometry) AND s1.wayid != s2.wayid;

DELETE FROM node WHERE nodeid IN
( SELECT mt1.nodeid FROM node mt1, node mt2 WHERE mt1.geometry ~=  mt2.geometry AND mt1.nodeid < mt2.nodeid );

UPDATE node
SET ways = subquery.ways
FROM ( SELECT node.nodeid, array_agg(way.wayid) AS ways FROM way, node WHERE ST_Intersects(way.geometry, node.geometry) GROUP BY node.nodeid ) AS subquery
WHERE subquery.nodeid = node.nodeid;

UPDATE way
SET nodes = subquery.nodes
FROM ( SELECT way.wayid, array_agg(node.nodeid) AS nodes FROM node, way WHERE ST_Intersects(way.geometry, node.geometry) GROUP BY way.wayid ) AS subquery
WHERE subquery.wayid = way.wayid;

UPDATE node
SET neighbors = subquery.neighbors
FROM ( SELECT node1.nodeid, array_agg(node2.nodeid) AS neighbors FROM node node1, node node2 WHERE node1.ways && node2.ways AND node1.nodeid != node2.nodeid GROUP BY node1.nodeid ) AS subquery
WHERE subquery.nodeid = node.nodeid;
