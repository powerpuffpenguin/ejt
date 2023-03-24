local env = import 'env.libsonnet';
{
  // get ext_strs
  user: std.extVar('USER'),
  // import lib from jpath
  verion: env.version,
  // import lib
  str: import './str.jsonnet',
  // import text
  text: importstr './str.jsonnet',
  // default ext_strs
  ejt: {
    dev: std.extVar('dev'),
    version: std.extVar('ejt.version'),
    os: std.extVar('ejt.os'),
    arch: std.extVar('ejt.arch'),
    go_version: std.extVar('ejt.go_version'),
    jsonnet: std.extVar('ejt.jsonnet'),
    dir: std.extVar('ejt.dir'),
  },
  // default func
  funcs: {
    readText: std.native('os.readText')(
      std.native('filepath.join')(std.extVar('ejt.dir'), 'ejt.jsonnet')
    ),
  },
}
