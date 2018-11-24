package stream

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
