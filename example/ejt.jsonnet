{
  version: 'v0.0.2',
  ext_strs: [  // std.extVar(key)
    'USER',  // get from environment
    'Generator=ejt',  // set <key>=<val>
  ],
  endpoints: [
    {
      output: './dst',  // redirect output structure to the directory.
      target: './envoy',  // target root directory.
      source: './src',  // source root directory.
      resources: [
        'static/envoy.jsonnet',
        'dynamic/envoy.jsonnet',
        'dynamic/cds.jsonnet',
        'dynamic/lds.jsonnet',
        'other/other.jsonnet',
      ],
      ext_strs: [
        'endpoint=envoy',
      ],
    },
  ],
}
