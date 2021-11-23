# Cache

Container to deploy a redis cache

Utility functions to interface with a redis cache

## Types

```
CacheDetails {
	Host        string
	IdleTimeout time.Duration
	MaxActive   int64 
	MaxIdle     int64 
	Port        int64
	Protocol    string
}
```

## Interfaces

```
CacheInterface {
    pool    <redis pool>
}

CacheInterface::Exec(args unknown[])->(unknown[])
```

## Containers

### Requirements

```
dnf install python3 golang podman podman-compose
```

### Configuration

```
./config/cache.json

{
    "container_port": 6379,
    "external_port": 3010,
    "max_samples": 5,
    "max_size_in_mb": "64mb"
}
```

### Scripts

```
python3 build.py

--config        config filepath
--templates     templates directory
--dest          destination directory
```

```
python3 run.py

--file          podman-compose filepath
```

```
python3 down.py

--file          podman-compose filepath
```


## License

BSD-3-Clause License
