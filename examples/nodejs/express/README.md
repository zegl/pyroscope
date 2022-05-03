# Pyroscope pull with static targets

This example demonstrates how Pyroscope can be used to scrape pprof profiles from remote nodejs targets.

### 1. Run Pyroscope server and demo application in docker containers

```shell
docker-compose up -d
```

Please note that we'd configured `pyroscope` to send data to server and tag it with region

```javascript
// Init pyroscope with the server name & region
Pyroscope.init({
  server: 'http://pyroscope:4040',
  tags: { region },
});
```

### 2. Observe profiling data

Profiling is more fun when the application does some work, so it shipped with built-in load generator.

Now that everything is set up, you can browse profiling data via [Pyroscope UI](http://localhost:4040).
