# redis-examples: Handle Delay Task With Redis

```shell
taskconsumer> go run main.go 
2018/05/18 10:21:08 start consume...
2018/05/18 10:21:15 payload: test task payload
2018/05/18 10:21:19 payload: test task payload
2018/05/18 10:21:23 payload: test task payload
```

## producer

* producer 将任务放到redis 的sorted set 中
* member 是任务的id
* score 是任务执行的时间

## consumer

consumer 将每隔一段时间从同一个sorted set中读取(zrange delayTaskListKey 0 0 withscores)第一个任务，然后判断score(执行时间) 与当前时间的是否相等，相等则执行任务，然后从set中删除这个member。
