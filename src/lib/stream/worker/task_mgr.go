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
	tasks     map[int]*Task
	lk        sync.RWMutex
	NextIndex int
}

var taskManager TaskManager

func (s *TaskManager) Task(id int) *Task {
	s.lk.RLock()
	defer s.lk.RUnlock()
	return s.tasks[id]
}

func (s *TaskManager) InsertTask(task *Task) int {
	s.lk.Lock()
	defer s.lk.Unlock()
	s.tasks[s.NextIndex] = task
	s.NextIndex++
	return s.NextIndex - 1
}

func (s *TaskManager) RemoveTask(task *Task) error {
	s.lk.Lock()
	defer s.lk.Unlock()
	delete(s.tasks, task.TaskId)
	return nil
}
