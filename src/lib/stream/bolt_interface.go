package stream

type BoltABC interface {
	Init()
	//Unpack()
	NextTuple()
	//Pack()
}
