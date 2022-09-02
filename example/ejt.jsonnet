{
  version: 'v0.0.1',
  endpoints: [
    {
      output: './dst',  // redirect output structure to the directory.
      target: './envoy',  // target root directory.
      source: './src',  // source root directory.
      resources: [
        'static/envoy.jsonnet',
      ],
    },
  ],
}
