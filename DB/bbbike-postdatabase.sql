SET search_path = topology,public;

update quality set type = replace(type, ';', '');

SELECT topology.CreateTopology('path_topo', 4326);
SELECT topology.AddTopoGeometryColumn('path_topo', 'public', 'path', 'topo_geom', 'LINESTRING');
SELECT topology.AddTopoGeometryColumn('path_topo', 'public', 'cyclepath', 'topo_geom', 'LINESTRING');
SELECT topology.AddTopoGeometryColumn('path_topo', 'public', 'greenpath', 'topo_geom', 'LINESTRING');

UPDATE path SET topo_geom = topology.toTopoGeom(geometry, 'path_topo', 1, 0.0);

insert into network (type, name, geometry, wayid)SELECT r.type, r.name, e.geom, r.id
FROM path_topo.edge e,
     path_topo.relation rel,
     path r
WHERE e.edge_id = rel.element_id
  AND rel.topogeo_id = (r.topo_geom).id;

DELETE FROM network WHERE networkid IN ( SELECT mt1.networkid FROM network mt1, network mt2 WHERE mt1.geometry ~=  mt2.geometry AND mt1.networkid < mt2.networkid );

insert into node (geometry) select distinct points.point from (select st_endpoint(geometry) as point from network
union
select st_startpoint(geometry) as point from network) as points;

UPDATE node
SET networks = subquery.networks
FROM ( SELECT node.id, array_agg(network.networkid) AS networks FROM network, node WHERE ST_Intersects(network.geometry, node.geometry) GROUP BY node.id ) AS subquery
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

UPDATE network SET attributes = array_append(attributes, row('greenway',info.type, info.geometry)::attribute) from (select path.networkid, greenpath.type as type, st_intersection(greenpath.geometry, path.geometry) as geometry from network as path, greenpath where st_intersects(greenpath.geometry, path.geometry) and not geometryType(st_intersection(path.geometry, greenpath.geometry)) = 'POINT') as info where info.networkid = network.networkid;
UPDATE network SET attributes = array_append(attributes, row('quality',info.type, info.geometry)::attribute) from (select path.networkid, quality.type as type, st_intersection(quality.geometry, path.geometry) as geometry from network as path, quality where st_intersects(quality.geometry, path.geometry) and not geometryType(st_intersection(path.geometry, quality.geometry)) = 'POINT') as info where info.networkid = network.networkid;
UPDATE network SET attributes = array_append(attributes, row('cyclepath', info.type, info.geometry)::attribute) from (select path.networkid, cyclepath.type as type, st_intersection(cyclepath.geometry, path.geometry) as geometry from network as path, cyclepath where st_intersects(cyclepath.geometry, path.geometry) and not geometryType(st_intersection(path.geometry, cyclepath.geometry)) = 'POINT') as info where info.networkid = network.networkid;
UPDATE network SET attributes = array_append(attributes, row('unlit', 'NL' , info.geometry)::attribute) from (select path.networkid, st_intersection(unlit.geometry, path.geometry) as geometry from network as path, unlitpath as unlit where st_intersects(unlit.geometry, path.geometry) and not geometryType(st_intersection(path.geometry, unlit.geometry)) = 'POINT' ) as info where info.networkid = network.networkid;
UPDATE network SET attributes = array_append(attributes, row('trafficlight', info.type , info.geometry)::attribute) from (select path.networkid, trafficlight.type, st_intersection(trafficlight.geometry, path.geometry) as geometry from network as path, trafficlight where st_intersects(trafficlight.geometry, path.geometry) and not trafficlight.type = 'X' ) as info where info.networkid = network.networkid;