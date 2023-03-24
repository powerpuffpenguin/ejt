{
  user: std.extVar('USER'),
  generator: std.extVar('Generator'),
  [if std.extVar('endpoint') == 'envoy' then 'envoy']: 'yes',
}
