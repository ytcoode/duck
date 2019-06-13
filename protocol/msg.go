package protocol

import "io"

type Msg struct {
	bs []byte
}

func NewMsg(msgID int) *Msg {
	m := &Msg{make([]byte, 0, 64)}
	if msgID > 0xff {
		panic("illegal msgID")
	}
	m.WriteByte(byte(msgID))
	return m
}

func (m *Msg) WriteTo(w io.Writer) error {
	l := len(m.bs)
	if l >= 0xffff {
		panic("Too big msg")
	}

	bs := make([]byte, 2, l+2)
	bs[0] = byte(l >> 8)
	bs[1] = byte(l)
	bs = append(bs, m.bs...)

	for len(bs) > 0 {
		n, err := w.Write(bs)
		if err != nil {
			return err
		}
		bs = bs[n:]
	}
	return nil
}

func ReadMsgFrom(r io.Reader) (*Msg, error) {
	bs := make([]byte, 2)
	_, err := io.ReadFull(r, bs)
	if err != nil {
		return nil, err
	}

	len := int(bs[0])<<8 | int(bs[1])
	m := &Msg{make([]byte, len)}

	_, err = io.ReadFull(r, m.bs)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Msg) WriteByte(b byte) {
	m.bs = append(m.bs, b)
}

func (m *Msg) WriteUint16(v uint16) {
	m.WriteByte(byte(v >> 8))
	m.WriteByte(byte(v))
}

func (m *Msg) WriteBytes(bs []byte) {
	if len(bs) > 0xffff {
		panic("WriteBytes")
	}
	m.WriteUint16(uint16(len(bs)))
	m.bs = append(m.bs, bs...)
}

func (m *Msg) ReadByte() byte {
	b := m.bs[0]
	m.bs = m.bs[1:]
	return b
}

func (m *Msg) ReadUint16() uint16 {
	b1, b2 := m.bs[0], m.bs[1]
	m.bs = m.bs[2:]
	return uint16(b1)<<8 | uint16(b2)
}

func (m *Msg) ReadBytes() []byte {
	len := int(m.ReadUint16())
	bs := m.bs[:len]
	m.bs = m.bs[len:]
	return bs
}

func (m *Msg) ReadMsgID() int {
	return int(m.ReadByte())
}
