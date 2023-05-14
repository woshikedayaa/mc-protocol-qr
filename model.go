package mcprotocolqr

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"net"
	"sync"
	"time"
)

// 4-byte request ID, 4-byte message type, 2-byte terminator.
const (
	headerSize int32 = 10

	TypeCommand = 2
	TypeAuth    = 3
)

type rconMessage struct {
	content     string
	contentType int32
	length      int32
	//这里用上级的指针指着
	id *int32
}

type RconServer struct {
	//基础数据
	addr    string
	port    string
	logined bool
	lastid  int32
	passwd  string

	connetcion *net.Conn
	mutex      *sync.RWMutex
	timeout    time.Duration
	//请求数据和返回数据
	BeginData rconMessage
	Result    rconMessage
}

func (r *rconMessage) DecodeMsg(b []byte) error {
	reader := bytes.NewReader(b)
	var (
		Length      int32
		ID          int32
		contentType int32
		content     []byte = make([]byte, 0)
	)

	if err := binary.Read(reader, binary.LittleEndian, &Length); err != nil {
		r.content = ""
		return err
	}
	r.length = Length

	if err := binary.Read(reader, binary.LittleEndian, &ID); err != nil {
		r.content = ""
		return err
	}
	r.id = &ID

	if err := binary.Read(reader, binary.LittleEndian, &contentType); err != nil {
		r.content = ""
		return err
	}
	r.contentType = contentType

	if err := binary.Read(reader, binary.LittleEndian, content); err != nil {
		r.content = ""
		return err
	}
	r.content = hex.EncodeToString(content)

	return nil
}

func (r *rconMessage) EncodeMsg() ([]byte, error) {
	buf := new(bytes.Buffer)
	for _, v := range []interface{}{
		r.length,
		r.id,
		r.contentType,
		[]byte(r.content),
		[]byte{0, 0}, // 2-byte terminator.
	} {
		if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func NewRconServer(addr string, port string, passwd string, timeout time.Duration) (*RconServer, error) {
	if timeout == 0 {
		timeout = 5 * time.Second
	}
	IP, err := net.ResolveTCPAddr("tcp", addr+":"+port)
	if err != nil {
		return nil, err
	}
	//初始化
	rs := new(RconServer)
	rs.addr = IP.String()
	rs.port = port
	rs.lastid = int32(0)
	rs.BeginData.id = &rs.lastid
	rs.mutex = new(sync.RWMutex)
	rs.passwd = passwd
	//憋说没用 只是显式声明下（
	//rs.logined=false

	if err := rs.connect(); err != nil {
		return nil, err
	}
	//登录 登录成功就把logined设置成true
	if err := rs.Login(); err != nil {
		return nil, err
	}
	rs.logined = true
	return rs, nil
}
