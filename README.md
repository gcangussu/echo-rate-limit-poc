# echo-rate-limit-poc

Test GCRA rate limiting for Echo server

The server is configured to give a rate limit of 10 req/s for each
client id. The client id is set with an `id` header on the HTTP
request, it can be any string.

1. Run Redis on port 6379. For example with docker:

```shell
docker run --rm -p 6379:6379 redis:6-alpine
```

2. Start the server and pipe data to the analyse.py script:

```shell
bash start_server.bash
```

3. Once the server is litening on port 1323, start the clients:

```shell
bash start_clients.bash
```

4. Stop the server with Ctrl + C to see the results. E.g.:

```
--- A ---
Total requests: 165
Total time (s): 14.749893
All requests rate (req/s): 11.186521827649868
2xx requests rate (req/s): 10.644145011763815

--- B ---
Total requests: 300
Total time (s): 14.951941
All requests rate (req/s): 20.064284630336623
2xx requests rate (req/s): 10.634070854078411
```

## Note

If you reduce the duration which the clients fire the requests
you may see the bursting ability provided by the fact that the
"bucket" for the client is initially empty.

Also, "2xx requets rate" is slightly higher than the 10 req/s
limit set on the server because of the initial requests "filling"
an empty "bucket". For example, with the limit of 10 req/s and a
client firing at 20 req/s, the client can get confortably fire 18
requests before being rate limited. The table below illustrates
this.

| req # | time (s) | done (req) | leaked (req) | filled (req) |
| ----- | -------- | ---------- | ------------ | ------------ |
| -     | -        | 0          | 0            | 0            |
| 1     | 0        | 1          | 0            | 1            |
| 2     | 0.05     | 2          | 0.5          | 1.5          |
| 3     | 0.1      | 3          | 1            | 2            |
| 4     | 0.15     | 4          | 1.5          | 2.5          |
| 5     | 0.2      | 5          | 2            | 3            |
| 6     | 0.25     | 6          | 2.5          | 3.5          |
| 7     | 0.3      | 7          | 3            | 4            |
| 8     | 0.35     | 8          | 3.5          | 4.5          |
| 9     | 0.4      | 9          | 4            | 5            |
| 10    | 0.45     | 10         | 4.5          | 5.5          |
| 11    | 0.5      | 11         | 5            | 6            |
| 12    | 0.55     | 12         | 5.5          | 6.5          |
| 13    | 0.6      | 13         | 6            | 7            |
| 14    | 0.65     | 14         | 6.5          | 7.5          |
| 15    | 0.7      | 15         | 7            | 8            |
| 16    | 0.75     | 16         | 7.5          | 8.5          |
| 17    | 0.8      | 17         | 8            | 9            |
| 18    | 0.85     | 18         | 8.5          | 9.5          |
| 19    | 0.9      | 19         | 9            | 10 (!)       |
| 20    | 0.95     | 20         | 9.5          | 10.5 (!!)    |
