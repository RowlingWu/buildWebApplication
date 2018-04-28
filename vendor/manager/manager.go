package manager

import (
	"errors"
	"os"
	"os/exec"
	"sync"
	"time"
)

type Task struct {
	ID       string   `json:"id"`
	Cmd      string   `json:"cmd"`
	Args     []string `json:"args"`
	Interval int      `json:"interval"`
}

// manager 全局的task管理器，负责创建定时任务、删除定时任务
var MyManager = &Manager{timer: make(map[string]*time.Timer), task: make(map[string]*Task)}

type Manager struct {
	lock  sync.Mutex
	timer map[string]*time.Timer
	task  map[string]*Task
}

// create 创建一个新的定时任务
func (manager *Manager) Create(task *Task) (bool, error) { // false代表该任务已存在。error返回其他类型错误
	if _, ok := manager.task[task.ID]; ok { // 如果已启动的任务里有该task.ID, 该task不应被启动
		return false, nil
	}
	manager.lock.Lock()
	defer manager.lock.Unlock()

	manager.task[task.ID] = task
	pipeline := func() { // pipline 执行task的命令
		exe := exec.Command(task.Cmd, task.Args...)
		exe.Stdout = os.Stdout
		exe.Start()
	}
	var doAfterFunc func()
	doAfterFunc = func() { // doAfterFunc 开始一个定时器
		pipeline()
		manager.timer[task.ID] = time.AfterFunc(time.Duration(task.Interval/1000)*time.Second, doAfterFunc)
	}
	manager.timer[task.ID] = time.AfterFunc(time.Duration(task.Interval/1000)*time.Second, doAfterFunc)

	return true, nil
}

// destroy 停止对task命令的执行
func (manager *Manager) Destroy(id string) error {
	if _, ok := manager.timer[id]; !ok {
		return errors.New("no such task id")
	}
	manager.lock.Lock()
	defer manager.lock.Unlock()

	delete(manager.task, id)
	manager.timer[id].Stop()
	delete(manager.timer, id)
	return nil
}
