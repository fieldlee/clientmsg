package utils

import (
	"fmt"
	"testing"
)

func TestGenerateKeyPair(t *testing.T) {
	data := []byte("123467890")
	i:=0
	for i < 10000 {
		pr,pb:=GenerateKeyPair(data,32)
		fmt.Println(PrivateKeyToBytes(pr))
		fmt.Println(PublicKeyToBytes(pb))
		i++
	}

}
