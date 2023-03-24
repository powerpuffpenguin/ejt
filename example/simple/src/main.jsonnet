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
}
