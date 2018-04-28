# mini GoCron

这是一个可以用HTTP协议控制的任务定时器。每隔一段时间，就会执行一条指定的系统命令。

## 用法

下载并安装：
`go get github.com/RowlingWu/ServiceComputing/minicron`

启动minicron:
`minicron [port]   // port是你想监听的端口，不输入时默认为8080`

例如想监听4567端口：
`minicron 4567`

## 创建任务

创建新的定时任务`date -R`，每隔3秒执行一次：
```
curl -d '
{
    "id":"minicron",
    "cmd":"date",
    "args":["-R"],
    "interval":3000
}
'  localhost:4567
```

若成功创建，返回 HTTP 200：
```
{
    "ok":true,
    "id":"minicron"
}
```

若任务已存在（该id已存在），返回 HTTP 409：
```
{
    "ok":false,
    "error":"The task minicron already exists."
}
```

## 终止任务

终止一个定时任务：`curl -X DELETE localhost:4567/{id}`

如要删除ID为minicron的任务：`curl -X DELETE localhost:4567/minicron`

若成功删除，返回 HTTP 200：
```
{
    "ok":true,
    "id":"minicron"
}
```

若没有该任务，返回 HTTP 404：
```
{
    ”ok”: false,
    ”error”: ”The task minicron is not found.”
}
```