package shared

type CollectorABC interface {
	Emit([]byte)
	//Ack()
	//Fail()
}

type SpoutABC interface {
	Init() error
	NextTuple(CollectorABC)
}

type BoltABC interface {
	Init() error
	//Unpack()
	Execute([]byte, CollectorABC)
	//Pack()
	//CollectorABC
}

type SinkABC interface {
	BoltABC
	CheckPoint()
}

//type Collector struct {
//	state
//}
//
//func (s *Collector) Emit() {
//	return
//}

//func (s *Collector) Ack() {
//	return
//}
//
//func (s *Collector) Fail() {
//	return
//}
