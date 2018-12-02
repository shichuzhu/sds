package shared

type CollectorABC interface {
	Emit([]byte)
	IssueCheckPoint()
	IssueStop()
	// Ack and Fail in the original storm design
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
