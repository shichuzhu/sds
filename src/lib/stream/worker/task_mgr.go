package worker

import (
	"fa18cs425mp/src/pb"
	"sync"
)

func GetTMgr() *TaskManager {
	return &taskManager
}

func IdFromCfg(cfg *pb.TaskCfg) int {
	return int(cfg.TaskId)
}

type TaskManager struct {
	tasks map[int]*Task
	lk    sync.RWMutex
}

var taskManager TaskManager

func (s *TaskManager) Task(id int) *Task {
	s.lk.RLock()
	defer s.lk.RUnlock()
	return s.tasks[id]
}

func (s *TaskManager) InsertTask(id int, task *Task) {
	s.lk.Lock()
	defer s.lk.Unlock()
	s.tasks[id] = task
}
