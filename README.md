# ejt

Small tool to convert jsonnet to yaml/json

English [中文](README.zh-Hant.md)

# Background

I want to use envoy to replace nginx in some small environments. These environments are not suitable for installing complex envoy control panels because these environments are simple or lightweight, but also need to be flexible enough to dynamically change settings at any time. envoy provides dynamic-resources-from-filesystem which is very suitable, but there are many problems with writing configuration directly using yaml. For example, it is impossible to disperse complex configurations into a single small file and import using techniques such as import/include. Most similar listeners/clusters repeat the same attribute values every time, and envoy will only reload dynamic when moving dynamic-resources specified content.

jsonnet can solve these yaml deficiencies very well, so I wrote this small tool to convert jsonnet to yaml and automatically copy yaml to the target path monitored by envoy dynamic-resources.

In addition, because the structure of jsonnet/json/yaml is consistent, the conversion from jsonnet to json is also supported by the way.

# How To Use

Create the **ejt.jsonnet** definition file, you can execute the following command to create the definition file in the current working directory:

```
ejt init
```

**ejt.jsonnet** defines where to get jsonnet from, where to output the transpiled code, and where to move or copy the output archives for an envoy-like monitoring system to trigger updates. which ends up looking like this:
```
{
  version: 'v0.0.1',
  endpoints: [
    {
      output: './dst',  // redirect output structure to the directory.
      target: './envoy',  // target root directory.
      source: './src',  // source root directory.
      resources: [
        'envoy.jsonnet',
      ],
    },
  ],
}
```

After writing your jsonnet, execute the following command to generate yaml and move to the target path:
```
ejt yaml -m
```