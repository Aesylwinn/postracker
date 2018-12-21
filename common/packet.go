package common

import (
    "bytes"
    "encoding/binary"
    "log"
    "time"
)

type PositionPayload struct {
    X float64
    Y float64
    Timestamp int64
}

func NewPositionPayload(x, y float64) PositionPayload {
    return PositionPayload{x, y, time.Now().UTC().Unix()}
}

func (pos *PositionPayload) Refresh() {
    pos.Timestamp = time.Now().UTC().Unix()
}

func (pos *PositionPayload) Encode() []byte {
    buffer := new(bytes.Buffer)
    binary.Write(buffer, binary.LittleEndian, pos)
    return buffer.Bytes()
}

func DecodePosition(data []byte) PositionPayload {
    var pos PositionPayload
    buffer := bytes.NewReader(data)
    err := binary.Read(buffer, binary.LittleEndian, &pos)
    if err != nil {
        log.Printf("Failed to decode position\n%v", err)
        return NewPositionPayload(0, 0)
    }
    return pos
}
