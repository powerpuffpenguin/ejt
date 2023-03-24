local def = import '../lib/def.libsonnet';
{
  admin: {
    access_log_path: '/dev/stdout',
    address: def.address('0.0.0.0', 8080),
  },
  node: {
    cluster: 'test-cluster',
    id: 'test-id',
  },
  dynamic_resources: {
    cds_config: {
      path: '/etc/envoy/cds.yaml',
    },
    lds_config: {
      path: '/etc/envoy/lds.yaml',
    },
  },
}
