# rdockerlogs

A utility for tailing the logs of remote servers that have Docker running with the syslog driver running.


### Usage

Tail a whole host:

`rdockerlogs -hosts=10.55.0.123 -format='{{index .Config.Labels "io.rancher.container.name"}}'`

Tail a multiple hosts:

`rdockerlogs -hosts=10.55.0.1,10.55.0.123 -format='{{index .Config.Labels "io.rancher.container.name"}}'`


### Things that don't work yet but I want to

Rancher flags like this which let's you tail specific Rancher containers:

`rdockerlogs -rancher=rancher.env -service='staging/mux' -format='{{index .Config.Labels "io.rancher.container.name"}}'`

Config files that let you save sets of command line flags:

`rdockerlogs -config=staging.yaml`

Other things to do:
- Specify more tail options like `-n` so you can see moar lines
- Allow setting the log file to tail in case the log file isn't `/var/log/syslog`
- Support for different formats of tagging in the log
- Document how to set up the syslog driver
