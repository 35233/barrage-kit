# barrage-forwarder

Configuring
----
barrage-forwarder is configured with a yaml file you specify:
```shell script
barrage-forwarder config.yaml
```

Here's a sample, with comments in-line to describe the settings:

```yaml
# The path of log file
log.path: 'barrage-forwarder.log'
# The list of source configurations
sources:
  - type: douyu
    roomids:
      - 196
      - 52004
# The list of output configurations
output:
  - type: file
    path: output.log
  - type: kafka
    brokers: 127.0.0.1:9092
    topic: test
  - type: elasticsearch
    bulkActions:  500
    bulkSize: 5000
    flushInterval: 600000
    urls:
      - http://127.0.0.1:9200
      - http://127.0.0.1:9210
```

Build
----
```shell script
go build github.com/35233/barrage-kit/barrage-forwarder
```