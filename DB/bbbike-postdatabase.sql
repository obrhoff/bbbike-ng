SET search_path = topology,public;

SELECT topology.CreateTopology('network_topo', 4326);
SELECT topology.AddTopoGeometryColumn('network_topo', 'public', 'streetway', 'topo_geom', 'LINESTRING');
SELECT topology.AddTopoGeometryColumn('network_topo', 'public', 'cycleway', 'topo_geom', 'LINESTRING');
SELECT topology.AddTopoGeometryColumn('network_topo', 'public', 'greenway', 'topo_geom', 'LINESTRING');

UPDATE streetway SET topo_geom = topology.toTopoGeom(geometry, 'network_topo', 1, 0.0);
UPDATE cycleway SET topo_geom = topology.toTopoGeom(geometry, 'network_topo', 2, 0.0);
UPDATE greenway SET topo_geom = topology.toTopoGeom(geometry, 'network_topo', 3, 0.0);

insert into network(foreignid, geometry, type) SELECT w.streetwayid, e.geom, w.type
FROM network_topo.edge e,
     network_topo.relation rel,
     streetway w
WHERE e.edge_id = rel.element_id
AND rel.topogeo_id = 1;

insert into network(foreignid, geometry, type) SELECT w.cycleid, e.geom, w.type
FROM network_topo.edge e,
     network_topo.relation rel,
     cycleway w
WHERE e.edge_id = rel.element_id
AND rel.topogeo_id = 2;

insert into network(foreignid, geometry, type) SELECT w.greenwayid, e.geom, w.type
FROM network_topo.edge e,
     network_topo.relation rel,
     greenway w
WHERE e.edge_id = rel.element_id
AND rel.topogeo_id = 3;

INSERT INTO node (geometry)
SELECT DISTINCT (ST_DumpPoints(ST_Intersection(s1.geometry, s2.geometry))).geom as geometry
FROM network s1, network s2
WHERE ST_Intersects(s2.geometry, s1.geometry) AND NOT s1.geometry = s2.geometry;

UPDATE node
SET networks = subquery.networks
FROM ( SELECT node.nodeid, array_agg(network.networkid) AS networks FROM network, node WHERE ST_Intersects(network.geometry, node.geometry) GROUP BY node.nodeid ) AS subquery
WHERE subquery.nodeid = node.nodeid;

UPDATE network
SET nodes = subquery.nodes
FROM ( SELECT network.networkid, array_agg(node.nodeid) AS nodes FROM node, network WHERE ST_Intersects(network.geometry, node.geometry) GROUP BY network.networkid ) AS subquery
WHERE subquery.networkid = network.networkid;

UPDATE node
SET neighbors = subquery.neighbors
FROM ( SELECT node1.nodeid, array_agg(node2.nodeid) AS neighbors FROM node node1, node node2 WHERE node1.networks && node2.networks AND node1.nodeid != node2.nodeid GROUP BY node1.nodeid ) AS subquery
WHERE subquery.nodeid = node.nodeid;

UPDATE node SET walkable = true WHERE array_length(neighbors, 1) > 1;
UPDATE node SET walkable = false WHERE array_length(neighbors, 1) < 2;