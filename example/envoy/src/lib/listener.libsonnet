local def = import './def.libsonnet';
{
  listener(address, port, filters): {
    address: def.address(address, port),
    filter_chains: filters,
  },
  http_connection_manager: {
    '@type': 'type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager',
    codec_type: 'AUTO',
    stat_prefix: 'ingress_http',
    access_log: [
      {
        name: 'envoy.access_loggers.stdout',
        typed_config: {
          '@type': 'type.googleapis.com/envoy.extensions.access_loggers.stream.v3.StdoutAccessLog',
        },
      },
    ],
    http_filters: [
      def.filters_http_router,
    ],
  },
  route_config(name, hosts): {
    name: name,
    virtual_hosts: hosts,
  },
  host(name, domains, routes): {
    name: name,
    domains: domains,
    routes: routes,
  },
  tls(certificate, key): {
    name: 'envoy.transport_sockets.tls',
    typed_config: {
      '@type': 'type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.DownstreamTlsContext',
      common_tls_context: {
        tls_certificates: [
          {
            certificate_chain: {
              inline_string: certificate,
            },
            private_key: {
              inline_string: key,
            },
          },
        ],
      },
    },
  },
}
