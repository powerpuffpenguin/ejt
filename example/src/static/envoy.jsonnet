local cluster = import '../lib/cluster.libsonnet';
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
  admin: {
    access_log_path: '/dev/stdout',
    address: def.address('0.0.0.0', 8080),
  },
  static_resources: {
    listeners: [
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
      ),
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
      ),
    ],
    clusters: [
      cluster.h2 {
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
  },
}
