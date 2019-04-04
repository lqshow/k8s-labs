## Job
Job 负责批处理任务，即仅执行一次的任务，它保证批处理任务的一个或多个Pod成功结束。

### Job spec

| spec                  | desc                                                |
| --------------------- | --------------------------------------------------- |
| RestartPolicy         | 仅支持 Never 或 OnFailure                           |
| completions           | 标志Job结束需要成功运行的Pod个数，默认为1           |
| parallelism           | 标志并行运行的Pod的个数，默认为1                    |
| activeDeadlineSeconds | 标志失败Pod的重试最大时间，超过这个时间不会继续重试 |

## CronJob

CronJob是基于调度的Job执行将会自动产生多个job，调度格式参考Linux的cron系统。