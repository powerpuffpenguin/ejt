local cluster = import '../lib/cluster.libsonnet';
local def = import '../lib/def.libsonnet';
{
  resources: [
    cluster.h2 {
      '@type': 'type.googleapis.com/envoy.config.cluster.v3.Cluster',
      name: 'google',
      type: 'LOGICAL_DNS',
      load_assignment: cluster.load_assignment(
        'google',
        [
          cluster.endpoint('www.google.com', 443),
        ]
      ),
    },
    cluster.h2 {
      '@type': 'type.googleapis.com/envoy.config.cluster.v3.Cluster',
      name: 'bing',
      type: 'LOGICAL_DNS',
      load_assignment: cluster.load_assignment(
        'bing',
        [
          cluster.endpoint('www.bing.com', 443),
        ]
      ),
    },
  ],
}
