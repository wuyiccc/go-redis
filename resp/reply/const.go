package reply

// pong
type PongReply struct {
}

var pongBytes = []byte("+PONG\r\n")

func MakePongReply() *PongReply {
	return &PongReply{}
}

func (r PongReply) ToBytes() []byte {
	return pongBytes
}

// ok
type OkReply struct{}

var okBytes = []byte("+OK\r\n")

var theOkReply = new(OkReply)

func MakeOkReply() *OkReply {
	return theOkReply
}

func (r *OkReply) ToBytes() []byte {
	return okBytes
}

// 空回复
type NullBulkReply struct {
}

var nullBulkBytes = []byte("$-1\r\n")
var theNullBulkReply = new(NullBulkReply)

func MakeNullBulkReply() *NullBulkReply {
	return theNullBulkReply
}

func (n NullBulkReply) ToBytes() []byte {
	return nullBulkBytes
}

// 空数组
type EmptyMultiBulkReply struct {
}

var emptyMultiBulkBytes = []byte("*0\r\n")

func (r *EmptyMultiBulkReply) ToBytes() []byte {
	return emptyMultiBulkBytes
}

// NoReply
type NoReply struct {
}

var noBytes = []byte("")

func (r *NoReply) ToBytes() []byte {
	return noBytes
}
