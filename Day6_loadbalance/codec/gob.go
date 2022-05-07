package codec

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
)

//conn 是由构建函数传入，通常是通过 TCP 或者 Unix 建立 socket 时得到的链接实例，
////dec 和 enc 对应 gob 的 Decoder 和 Encoder，buf 是为了防止阻塞而创建的带缓冲的 Writer，一般这么做能提升性能

type GobCodec struct {
	conn io.ReadWriteCloser //载体
	buf  *bufio.Writer
	dec  *gob.Decoder //反编码
	enc  *gob.Encoder //gob编码
}

//这是用来检测GobCodec是否实现了Codec接口，如果没有实现该接口则编译的时候就会报错
var _Codec = (*GobCodec)(nil)

func NewGobCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	return &GobCodec{
		conn: conn,
		buf:  buf,
		dec:  gob.NewDecoder(conn),
		enc:  gob.NewEncoder(conn),
	}
}

//接着实现 ReadHeader、ReadBody、Write 和 Close 方法

func (c *GobCodec) ReadHeader(h *Header) error {
	return c.dec.Decode(h)
}

func (c *GobCodec) ReadBody(body interface{}) error {
	return c.dec.Decode(body)
}

func (c *GobCodec) Write(h *Header, body interface{}) (err error) {
	defer func() {
		_ = c.buf.Flush()
		if err != nil {
			_ = c.Close()
		}
	}()

	if err := c.enc.Encode(h); err != nil {
		log.Println("rpc codec:gob error encoding header:", err.Error())
		return err
	}

	if err := c.enc.Encode(body); err != nil {
		log.Println("rpc codec:gob error encoding body:", err.Error())
		return err
	}
	return nil
}

func (c *GobCodec) Close() error {
	return c.conn.Close()
}
