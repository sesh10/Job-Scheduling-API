# Job Scheduling API
Created job scheduling API server using two ways.
- SuperAlign_Assgn: Used simple thread safe map for smaller scale.
- Redis_SuperAlign: Used redis for in-memory storage


## Steps to Run - SuperAlign_Assgn
Open a terminal/shell instance. Navigate inside SuperAlign_Assgn directory
```
cd SuperAlign_Assgn
```

Now, run the following command to start the server:
```
go run main.go
```

To access the endpoints you can open another terminal/shell and run the following commands:
- To create a job:
```
curl -X POST http://localhost:8080/jobs
```

- To fetch unique job and status:
```
curl http://localhost:8080/jobs/<job_id>
```


## Steps to Run - Redis_SuperAlign
To run Redis, you can do this in two ways:
- Using Docker (Assuming Ubuntu system)
  - Install Docker from the following link: https://docs.docker.com/engine/install/ubuntu/
  - Start docker engine from terminal
  - Use following command to run Redis image
    ```
    docker run --name redis-server -p 6379:6379 redis
    ```

- Local Redis Installation
  - Install redis from the following link: https://redis.io/docs/latest/operate/oss_and_stack/install/archive/install-redis/
  - Once Redis is up and running try the following command to test connection.
    ```
    redis-cli
    ```
  - If you get "PONG" as the response you are good to go


Open a terminal/shell instance. Navigate inside SuperAlign_Assgn directory
```
cd Redis_SuperAlign
```

Now, run the following command to start the server:
```
go run main.go
```

To access the endpoints you can open another terminal/shell instance and run the following commands:
- To create a job:
```
curl -X POST http://localhost:8080/jobs
```

- To fetch unique job and status:
```
curl http://localhost:8080/jobs/<job_id>
```

- To check job status in Redis Database, open another terminal/shell instance and run:
```
redis-cli GET job:<job_id>
```
