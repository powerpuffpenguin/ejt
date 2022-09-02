{
  address(address, port): {
    socket_address: {
      address: address,
      port_value: port,
    },
  },
  filters_http_router: {
    name: 'envoy.filters.http.router',
    typed_config: {
      '@type': 'type.googleapis.com/envoy.extensions.filters.http.router.v3.Router',
    },
  },
}
