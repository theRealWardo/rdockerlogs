rdockerlogs
-----------

- Tail a whole host: `rdockerlogs -hosts=10.55.0.123 -format='{{index .Config.Labels "io.rancher.container.name"}}'`
- Tail a multiple hosts: `rdockerlogs -hosts=10.55.0.1,10.55.0.123 -format='{{index .Config.Labels "io.rancher.container.name"}}'`
- Tail specific Rancher containers: `rdockerlogs -rancher=rancher.env -service='staging/mux' -format='{{index .Config.Labels "io.rancher.container.name"}}'`
- Use a config: `rdockerlogs -config=staging.yaml`

Things that don't work yet:
- Rancher flags
- Config file

Other things to do:
- Add ability to specify more tail options
