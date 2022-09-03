local def = import '../lib/def.libsonnet';
local listener = import '../lib/listener.libsonnet';
local routes = [
  {
    match: {
      prefix: '/',
      headers: [
        {
          name: ':authority',
          string_match: { exact: 'google.com' },
        },
      ],
    },
    route: {
      cluster: 'google',
      host_rewrite_literal: 'www.google.com',
    },
  },
  {
    match: {
      prefix: '/',
      headers: [
        {
          name: ':authority',
          string_match: { exact: 'bing.com' },
        },
      ],
    },
    route: {
      cluster: 'bing',
      host_rewrite_literal: 'www.bing.com',
    },
  },
];
{
  resources: [
    listener.listener(
      '0.0.0.0', 80, [{
        filters: [
          {
            name: 'envoy.filters.network.http_connection_manager',
            typed_config: listener.http_connection_manager {
              route_config: listener.route_config('search_route', [
                listener.host('service', ['*'], routes),
              ]),
            },
          },
        ],
      }],
    ) {
      name: 'listener_http',
      '@type': 'type.googleapis.com/envoy.config.listener.v3.Listener',
    },
    listener.listener(
      '0.0.0.0', 443, [{
        filters: [
          {
            name: 'envoy.filters.network.http_connection_manager',
            typed_config: listener.http_connection_manager {
              stat_prefix: 'ingress_https',
              route_config: listener.route_config('search_route', [
                listener.host('service', ['*'], routes),
              ]),
            },
          },
        ],
        transport_socket: listener.tls(
          importstr '../lib/test.pem',
          importstr '../lib/test.key',
        ),
      }],
    ) {
      name: 'listener_https',
      '@type': 'type.googleapis.com/envoy.config.listener.v3.Listener',
    },
  ],
}
