# HTTP Server

## Timeout testing

```sh
telnet localhost 8080
```

Check that the connection times out after set interval.  

## Test

```sh
hey -z 3s -c 50 http://localhost:8080
```

Results

```sh
Summary:
  Total:	3.0047 secs
  Slowest:	0.0360 secs
  Fastest:	0.0000 secs
  Average:	0.0018 secs
  Requests/sec:	27576.8085
  

Response time histogram:
  0.000 [1]	|
  0.004 [73085]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.007 [8376]	|■■■■■
  0.011 [1142]	|■
  0.014 [194]	|
  0.018 [23]	|
  0.022 [18]	|
  0.025 [10]	|
  0.029 [7]	|
  0.032 [2]	|
  0.036 [3]	|


Latency distribution:
  10% in 0.0002 secs
  25% in 0.0007 secs
  50% in 0.0013 secs
  75% in 0.0024 secs
  90% in 0.0039 secs
  95% in 0.0051 secs
  99% in 0.0084 secs

Details (average, fastest, slowest):
  DNS+dialup:	0.0000 secs, 0.0000 secs, 0.0360 secs
  DNS-lookup:	0.0000 secs, 0.0000 secs, 0.0140 secs
  req write:	0.0001 secs, 0.0000 secs, 0.0109 secs
  resp wait:	0.0010 secs, 0.0000 secs, 0.0220 secs
  resp read:	0.0005 secs, 0.0000 secs, 0.0151 secs

Status code distribution:
  [200]	82861 responses
```
