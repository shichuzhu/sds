package shared

type CollectorABC interface {
	Emit()
	Ack()
	Fail()
}

type SpoutABC interface {
	Init()
	NextTuple()
}

type BoltABC interface {
	Init()
	//Unpack()
	Execute()
	//Pack()
	CollectorABC
}

type SinkABC interface {
	BoltABC
	CheckPoint()
}

type Collector struct {
	// state
}

func (s *Collector) Emit() {
	return
}

func (s *Collector) Ack() {
	return
}

func (s *Collector) Fail() {
	return
}
