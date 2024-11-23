package utils

import (
	"fmt"
	"time"
)

type ITask interface {
	onTaskInit()
	runTask() bool
	canRemove() bool
	canRun() bool
	getName() string
}

type BaseTask struct {
	name     string                  //任务名称
	task     func(params []any) bool //执行的任务
	duration int                     //间隔 毫秒
	maxTimes int                     //最大执行次数 -1为无限
	params   []any                   //参数
	lastTime int64                   //上次执行时间
}

func (t *BaseTask) onTaskInit() {
}
func (t *BaseTask) runTask() bool {
	t.lastTime = GetNow()
	if t.maxTimes > 0 {
		t.maxTimes--
	} else if t.maxTimes == 0 {
		return false
	}
	return t.task(t.params)
}
func (t *BaseTask) canRemove() bool {
	return t.maxTimes == 0
}

func (t *BaseTask) canRun() bool {
	return GetNow()-t.lastTime >= int64(t.duration)
}

func (t *BaseTask) getName() string {
	return t.name
}

var taskList List[ITask]
var isInit bool = false

func addTask(task ITask) {
	checkInitTask()
	if task.canRun() {
		task.onTaskInit()
		taskList.Add(task)

	}
}

func RemoveTask(task ITask) {
	taskList.Remove(task)
}

func CreateTaskWithDuration(name string, duration int, task func(param []any) bool, params ...any) ITask {
	baseTask := BaseTask{
		name:     name,
		task:     task,
		duration: duration,
		maxTimes: -1,
		params:   params,
	}
	addTask(&baseTask)
	return &baseTask
}

func checkInitTask() {
	if !isInit {
		taskList = List[ITask]{}
		isInit = true
		go func() {
			for {
				taskList.ForEach(func(task ITask) {
					go func() {
						if task.canRemove() {
							taskList.Remove(task)
							return
						}
						if task.canRun() {
							if !task.runTask() {
								taskList.Remove(task)
							}
							if log.IsDebug() {
								log.Debug(fmt.Sprintf("task %s run", task.getName()))
							}
						}
					}()
				})
				time.Sleep(100 * time.Millisecond)
			}
		}()
	}
}
