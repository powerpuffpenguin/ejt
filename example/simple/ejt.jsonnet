{
  version: 'v0.0.2',
  ext_strs: [  // std.extVar(key)
    'USER',  // get from environment
    // set <key>=<val>
    'Generator=ejt',
  ],
  endpoints: [
    {
      output: './dst',  // redirect output structure to the directory.
      target: './target',  // target root directory.
      source: './src',  // source root directory.
      resources: [
        'main.jsonnet',
      ],
      ext_strs: [
        'endpoint=envoy',
      ],
      jpath: [
        'lib',
      ],
    },
  ],
}
