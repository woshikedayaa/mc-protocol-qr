package mcprotocolqr

import (
	"net"
	"time"
)

func (r *RconServer) Get() string {
	return r.Result.content
}

// 请求
func (r *RconServer) request(content string, contentType int32) error {
	msgraw := rconMessage{
		contentType: contentType,
		content:     content,
		length:      int32(len(content) + int(headerSize)),
		id:          &r.lastid,
	}
	var msgb = make([]byte, 0)
	var err error
	if msgb, err = msgraw.EncodeMsg(); err != nil {
		return err
	}
	//发送数据
	(*(r.connetcion)).SetWriteDeadline(time.Now().Add(r.timeout))
	_, err = (*(r.connetcion)).Write(msgb)

	if err != nil {
		return err
	}

	//接受数据
	(*(r.connetcion)).SetReadDeadline(time.Now().Add(r.timeout))
	_, err = (*(r.connetcion)).Read(msgb)
	if err != nil {
		return err
	}
	if err := r.Result.DecodeMsg(msgb); err != nil {
		return err
	}
	//id加一 方便下次使用
	r.lastid++
	return nil
}

// 用户正常输入命令
func (r *RconServer) Run(command string) error {
	//这里记录数据
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if err := r.request(command, TypeCommand); err != nil {
		return err
	}
	return nil
}

func (r *RconServer) Login() error {
	return r.request(r.passwd, TypeAuth)
}

func (r *RconServer) connect() error {
	//你要重新连接了 肯定是logout了
	if r.logined != false {
		return nil
	}
	conn, err := net.DialTimeout("tcp", r.addr, r.timeout)

	if err != nil {
		return err
	}
	r.connetcion = &conn
	return nil
}

func (r *RconServer) Close() error {
	r.logined = false

	if err := (*(r.connetcion)).Close(); err != nil {
		return err
	}
	return nil
}
