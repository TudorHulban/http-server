# HTTP Server

## Timeout testing

```sh
telnet localhost 8080
```

Check that the connection times out after set interval.  

## Test HTTP

```sh
hey -z 3s -c 50 http://localhost:8080
```

### Results

```sh
hey -z 3s -c 50 http://localhost:8080

Summary:
  Total:	3.0005 secs
  Slowest:	0.0098 secs
  Fastest:	0.0000 secs
  Average:	0.0002 secs
  Requests/sec:	267761.6470
  

Response time histogram:
  0.000 [1]	|
  0.001 [797780]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.002 [4645]	|
  0.003 [624]	|
  0.004 [237]	|
  0.005 [60]	|
  0.006 [12]	|
  0.007 [1]	|
  0.008 [0]	|
  0.009 [19]	|
  0.010 [31]	|


Latency distribution:
  10% in 0.0001 secs
  25% in 0.0001 secs
  50% in 0.0001 secs
  75% in 0.0002 secs
  90% in 0.0003 secs
  95% in 0.0005 secs
  99% in 0.0009 secs

Details (average, fastest, slowest):
  DNS+dialup:	0.0000 secs, 0.0000 secs, 0.0098 secs
  DNS-lookup:	0.0000 secs, 0.0000 secs, 0.0016 secs
  req write:	0.0000 secs, 0.0000 secs, 0.0097 secs
  resp wait:	0.0001 secs, 0.0000 secs, 0.0090 secs
  resp read:	0.0000 secs, 0.0000 secs, 0.0082 secs

Status code distribution:
  [200]	803410 responses
```

## Test HTTPS

```sh
curl -k -v https://localhost:443
```
