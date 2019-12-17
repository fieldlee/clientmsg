package utils

import (
	"bytes"
	"encoding/binary"
	"clientmsg/model"
	"fmt"
)

func JoinHeadAndBody(info model.HeadInfo,in []byte)[]byte{
	tagByte := StringToBytes(info.Tag)
	versionByte := Int16ToBytes(info.Version)
	clientTypeByte := Int16ToBytes(info.ClientType)
	headLengthByte := Int16ToBytes(info.HeadLength)
	CompressWayBYte := Uint8ToBytes(info.CompressWay)
	EncryptionBYte := Uint8ToBytes(info.Encryption)
	SigBYte := Uint8ToBytes(info.Sig)
	FormatBYte := Uint8ToBytes(info.Format)
	NetFlagBYte := Uint8ToBytes(info.NetFlag)
	Back1BYte := Uint8ToBytes(info.Back1)
	BufSizeBYte := Int32ToBytes(info.BufSize)
	UncompressedSizeByte := Int32ToBytes(info.UncompressedSize)
	Back2Byte := Int32ToBytes(info.Back2)
	return BytesJoin(tagByte,versionByte,clientTypeByte,headLengthByte,CompressWayBYte,EncryptionBYte,SigBYte,
		FormatBYte,NetFlagBYte,Back1BYte,BufSizeBYte,UncompressedSizeByte,Back2Byte,in)
}

func BytesJoin(blist ...[]byte)[]byte{
	bytesinfo := make([]byte,0)
	for _,b := range blist  {
		bytesinfo = append(bytesinfo,b...)
	}
	return bytesinfo
}

func BytesToString(b []byte)string{
	nb := make([]byte,0)
	for _,t := range b {
		x := fmt.Sprintf("%v",t)
		if x != "0"{
			nb = append(nb,t)
		}
	}
	return string(nb)
}

func StringToBytes(s string)[]byte{
	b := []byte(s)
	return b[:8]
}


func Int16ToBytes(n int16)[]byte{
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.LittleEndian, n)
	if err != nil {
		return nil
	}
	return buffer.Bytes()[:2]
}

func Uint8ToBytes(n uint8)[]byte{
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.LittleEndian, n)
	if err != nil {
		return nil
	}
	return buffer.Bytes()[:1]
}

func Uint32ToBytes(n uint32)[]byte{
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.LittleEndian, n)
	if err != nil {
		return nil
	}
	return buffer.Bytes()[:4]
}

func Int32ToBytes(n int32)[]byte{
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.LittleEndian, n)
	if err != nil {
		return nil
	}
	return buffer.Bytes()[:4]
}

//字节转换成整形
func BytesToInt16(b []byte) int16 {
	//b := ClearBytes(by)
	bytesBuffer := bytes.NewBuffer(b)
	var x int16
	binary.Read(bytesBuffer, binary.LittleEndian, &x)
	return x
}
// unsigned char -->  C.uchar -->  uint8
func BytesToUInt8(b []byte) uint8 {
	//b := ClearBytes(by)
	bytesBuffer := bytes.NewBuffer(b)
	var x uint8
	binary.Read(bytesBuffer, binary.LittleEndian, &x)
	return x
}

func BytesToInt32(b []byte) int32 {
	//b := ClearBytes(by)
	bytesBuffer := bytes.NewBuffer(b)
	var x int32
	binary.Read(bytesBuffer, binary.LittleEndian, &x)
	return int32(x)
}
