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

UPDATE network SET normal = array_cat(normal, done.array) from
(select final.networkid, array_agg(('CA', final.type, final.real)::attribute) as array from (
select next.networkid as networkid, next.type, st_linemerge(next.real) as real, next.normal from (select test.networkid, test.type, test.real, test.normal, geometryType(test.real) as geotype from (select networkid, cyclepath.type as type, st_intersection(cyclepath.geometry, network.geometry) as real, st_intersection(network.geometry, cyclepath.geometry) as normal from network, cyclepath where st_intersects(network.geometry, cyclepath.geometry)) as test) as next where (next.geotype = 'LINESTRING' OR next.geotype = 'MULTILINESTRING') and st_orderingequals(next.normal, next.real)) as final group by final.networkid) as done where done.networkid = network.networkid;

UPDATE network SET reversed = array_cat(reversed, done.array) from
(select final.networkid, array_agg(('CA', final.type, final.real)::attribute) as array from (
select next.networkid as networkid, next.type, st_linemerge(next.real) as real, next.normal from (select test.networkid, test.type, test.real, test.normal, geometryType(test.real) as geotype from (select networkid, cyclepath.type as type, st_intersection(cyclepath.geometry, network.geometry) as real, st_intersection(network.geometry, cyclepath.geometry) as normal from network, cyclepath where st_intersects(network.geometry, cyclepath.geometry)) as test) as next where (next.geotype = 'LINESTRING' OR next.geotype = 'MULTILINESTRING') and not st_orderingequals(next.normal, next.real)) as final group by final.networkid) as done where done.networkid = network.networkid;

UPDATE network SET normal = array_cat(normal, done.array) from
(select final.networkid, array_agg(('QA', final.type, final.real)::attribute) as array from (
select next.networkid as networkid, next.type, st_linemerge(next.real) as real, next.normal from (select test.networkid, test.type, test.real, test.normal, geometryType(test.real) as geotype from (select networkid, quality.type as type, st_intersection(quality.geometry, network.geometry) as real, st_intersection(network.geometry, quality.geometry) as normal from network, quality where st_intersects(network.geometry, quality.geometry)) as test) as next where (next.geotype = 'LINESTRING' OR next.geotype = 'MULTILINESTRING') and st_orderingequals(next.normal, next.real)) as final group by final.networkid) as done where done.networkid = network.networkid;

UPDATE network SET reversed = array_cat(reversed, done.array) from
(select final.networkid, array_agg(('QA', final.type, final.real)::attribute) as array from (
select next.networkid as networkid, next.type, st_linemerge(next.real) as real, next.normal from (select test.networkid, test.type, test.real, test.normal, geometryType(test.real) as geotype from (select networkid, quality.type as type, st_intersection(quality.geometry, network.geometry) as real, st_intersection(network.geometry, quality.geometry) as normal from network, quality where st_intersects(network.geometry, quality.geometry)) as test) as next where (next.geotype = 'LINESTRING' OR next.geotype = 'MULTILINESTRING') and not st_orderingequals(next.normal, next.real)) as final group by final.networkid) as done where done.networkid = network.networkid;

UPDATE network SET normal = array_cat(normal, done.array) from
(select final.networkid, array_agg(('UA', 'NL', final.real)::attribute) as array from (
select next.networkid as networkid, st_linemerge(next.real) as real, next.normal from (select test.networkid, test.real, test.normal, geometryType(test.real) as geotype from (select networkid, st_intersection(unlitpath.geometry, network.geometry) as real, st_intersection(network.geometry, unlitpath.geometry) as normal from network, unlitpath where st_intersects(network.geometry, unlitpath.geometry)) as test) as next where (next.geotype = 'LINESTRING' OR next.geotype = 'MULTILINESTRING') and st_orderingequals(next.normal, next.real)) as final group by final.networkid) as done where done.networkid = network.networkid;

UPDATE network SET reversed = array_cat(reversed, done.array) from
(select final.networkid, array_agg(('UA', 'NL', final.real)::attribute) as array from (
select next.networkid as networkid, st_linemerge(next.real) as real, next.normal from (select test.networkid, test.real, test.normal, geometryType(test.real) as geotype from (select networkid, st_intersection(unlitpath.geometry, network.geometry) as real, st_intersection(network.geometry, unlitpath.geometry) as normal from network, unlitpath where st_intersects(network.geometry, unlitpath.geometry)) as test) as next where (next.geotype = 'LINESTRING' OR next.geotype = 'MULTILINESTRING') and not st_orderingequals(next.normal, next.real)) as final group by final.networkid) as done where done.networkid = network.networkid;

UPDATE network SET normal = array_cat(normal, done.array) from
(select final.networkid, array_agg(('HA', final.type, final.real)::attribute) as array from (
select next.networkid as networkid, next.type, st_linemerge(next.real) as real, next.normal from (select test.networkid, test.type, test.real, test.normal, geometryType(test.real) as geotype from (select networkid, handicap.type as type, st_intersection(handicap.geometry, network.geometry) as real, st_intersection(network.geometry, handicap.geometry) as normal from network, handicap where st_intersects(network.geometry, handicap.geometry)) as test) as next where (next.geotype = 'LINESTRING' OR next.geotype = 'MULTILINESTRING') and st_orderingequals(next.normal, next.real)) as final group by final.networkid) as done where done.networkid = network.networkid;

UPDATE network SET reversed = array_cat(reversed, done.array) from
(select final.networkid, array_agg(('HA', final.type, final.real)::attribute) as array from (
select next.networkid as networkid, next.type, st_linemerge(next.real) as real, next.normal from (select test.networkid, test.type, test.real, test.normal, geometryType(test.real) as geotype from (select networkid, handicap.type as type, st_intersection(handicap.geometry, network.geometry) as real, st_intersection(network.geometry, handicap.geometry) as normal from network, handicap where st_intersects(network.geometry, handicap.geometry)) as test) as next where (next.geotype = 'LINESTRING' OR next.geotype = 'MULTILINESTRING') and not st_orderingequals(next.normal, next.real)) as final group by final.networkid) as done where done.networkid = network.networkid;

UPDATE network SET defaults = array_cat(defaults, done.array) from
(select final.networkid, array_agg(('GA', final.type, final.real)::attribute) as array from (
select next.networkid as networkid, next.type, st_linemerge(next.real) as real, next.normal from (select test.networkid, test.type, test.real, test.normal, geometryType(test.real) as geotype from (select networkid, greenpath.type as type, st_intersection(greenpath.geometry, network.geometry) as real, st_intersection(network.geometry, greenpath.geometry) as normal from network, greenpath where st_intersects(network.geometry, greenpath.geometry)) as test) as next where (next.geotype = 'LINESTRING' OR next.geotype = 'MULTILINESTRING') ) as final group by final.networkid) as done where done.networkid = network.networkid;

UPDATE network SET defaults = array_cat(defaults, done.array) from
(select final.networkid, array_agg(('GA', final.type, final.real)::attribute) as array from (
select next.networkid as networkid, next.type, next.real as real, next.normal from (select test.networkid, test.type, test.real, test.normal, geometryType(test.real) as geotype from (select networkid, trafficlight.type as type, st_intersection(trafficlight.geometry, network.geometry) as real, st_intersection(network.geometry, trafficlight.geometry) as normal from network, trafficlight where st_intersects(network.geometry, trafficlight.geometry)) as test) as next where (next.geotype = 'POINT') ) as final group by final.networkid) as done where done.networkid = network.networkid;TE network SET defaults = array_cat(defaults, attribute.array) from (select info.networkid, array_agg(('TA', info.type, info.geometry)::attribute) as array from (select path.networkid, trafficlight.type as type, st_intersection(trafficlight.geometry, path.geometry) as geometry from network as path, trafficlight where st_intersects(trafficlight.geometry, path.geometry) and not trafficlight.type = 'X')  as info group by info.networkid) as attribute where attribute.networkid = network.networkid;