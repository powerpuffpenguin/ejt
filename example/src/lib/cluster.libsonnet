local def = import './def.libsonnet';
{
  h2c: {
    connect_timeout: '1s',
    type: 'STRICT_DNS',
    lb_policy: 'ROUND_ROBIN',
  },
  h2: self.h2c {
    transport_socket: {
      name: 'envoy.transport_sockets.tls',
      typed_config: {
        '@type': 'type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.UpstreamTlsContext',
      },
    },
  },
  load_assignment(name, endpoints): {
    cluster_name: name,
    endpoints: endpoints,
  },
  endpoint(addr, port): {
    lb_endpoints: [
      {
        endpoint: {
          address: def.address(addr, port),
        },
      },
    ],
  },
}
