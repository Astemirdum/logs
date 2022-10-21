# REST LogApp на Go

## API:
- write log
- read log
- read logs

## The following concepts are applied in app:
- The Clean Architecture
- Nginx (as load-balancer of 2 app-instance)
- Postgres Replication - logical: replication factor=3
- Graceful Shutdown
- Running app in docker containers

#### run app
```
make run
```
#### stop app
```
make stop
```

#### to make requests 
```
cd calls
```