SET search_path = topology,public;

update quality set type = replace(type, ';', '');

SELECT topology.CreateTopology('path_topo', 4326);
SELECT topology.AddTopoGeometryColumn('path_topo', 'public', 'path', 'topo_geom', 'LINESTRING');
SELECT topology.AddTopoGeometryColumn('path_topo', 'public', 'cyclepath', 'topo_geom', 'LINESTRING');
SELECT topology.AddTopoGeometryColumn('path_topo', 'public', 'greenpath', 'topo_geom', 'LINESTRING');

UPDATE path SET topo_geom = topology.toTopoGeom(geometry, 'path_topo', 1, 0.0);
UPDATE cyclepath SET topo_geom = topology.toTopoGeom(geometry, 'path_topo', 2, 0.0);
UPDATE greenpath SET topo_geom = topology.toTopoGeom(geometry, 'path_topo', 3, 0.0);

insert into network (type, name, geometry, wayid)SELECT r.type, r.name, e.geom, r.id
FROM path_topo.edge e,
     path_topo.relation rel,
     path r
WHERE e.edge_id = rel.element_id
  AND rel.topogeo_id = (r.topo_geom).id;

insert into network (type, geometry)SELECT r.type, e.geom
FROM path_topo.edge e,
     path_topo.relation rel,
     greenpath r
WHERE e.edge_id = rel.element_id
  AND rel.topogeo_id = (r.topo_geom).id;

insert into network (type, geometry)SELECT r.type, e.geom
FROM path_topo.edge e,
     path_topo.relation rel,
     cyclepath r
WHERE e.edge_id = rel.element_id
  AND rel.topogeo_id = (r.topo_geom).id;



insert into node (geometry) select distinct points.point from (select st_endpoint(geometry) as point from network
union
select st_startpoint(geometry) as point from network) as points;

UPDATE node
SET networks = subquery.networks
FROM ( SELECT node.id, array_agg(network.id) AS networks FROM network, node WHERE ST_Intersects(network.geometry, node.geometry) GROUP BY node.id ) AS subquery
WHERE subquery.id = node.id;

UPDATE network
SET nodes = subquery.nodes
FROM (select networks.networkid as networkid, array_agg(n1.id) as nodes from node as n1, (select networkid from network) as networks where n1.networks @> ARRAY[networkid] group by networkid) as subquery where network.networkid = subquery.networkid;

UPDATE node
SET neighbors = subquery.neighbors
FROM ( SELECT node1.id, array_agg(node2.id) AS neighbors FROM node node1, node node2 WHERE node1.networks && node2.networks AND node1.id != node2.id GROUP BY node1.id ) AS subquery
WHERE subquery.id = node.id;

UPDATE node SET walkable = true WHERE array_length(neighbors, 1) > 1;
UPDATE node SET walkable = false WHERE array_length(neighbors, 1) < 2;

update network set attributes = hstore('quality',quality.type) from (select network.networkid, quality.type from network, quality where st_intersects(quality.geometry, network.geometry) and not geometryType(st_intersection(network.geometry, quality.geometry)) = 'POINT') as quality where  quality.networkid = network.networkid;
UPDATE network SET attributes = attributes || hstore('greenway',subquery.type) from (select networkid, greenpath.type, geometryType(st_intersection(path.geometry, greenpath.geometry)) from network as path, greenpath where st_intersects(greenpath.geometry, path.geometry) and not geometryType(st_intersection(path.geometry, greenpath.geometry)) = 'POINT') as subquery where network.networkid = subquery.networkid;
UPDATE network SET attributes = attributes || hstore('unlit', 'NL') from (select networkid, geometryType(st_intersection(path.geometry, unlitpath.geometry)) from network as path, unlitpath where st_intersects(unlitpath.geometry, path.geometry) and not geometryType(st_intersection(path.geometry, unlitpath.geometry)) = 'POINT') as subquery where network.networkid = subquery.networkid;
update node set trafficlight = true from (SELECT node.id as resultid  FROM node RIGHT JOIN trafficlight ON st_equals(node.geometry, trafficlight.geometry)) as result where result.resultid = node.id;
update node set trafficlight = false from  (select node.id as resultid from node where trafficlight is null) as result where result.resultid = nodeid;