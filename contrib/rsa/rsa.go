package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	//	"crypto/sha1"
	"crypto/x509"

	//	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"

	"math/big"
	"strconv"
)

func fromBase10(base10 string) *big.Int {
	i, ok := new(big.Int).SetString(base10, 10)
	if !ok {
		panic("bad number: " + base10)
	}
	return i
}

func parseInt16(n string) int {
	i, _ := strconv.ParseInt(n, 16, 64)
	return int(i)
}
func FromBase16(base16 string) *big.Int {
	i, ok := new(big.Int).SetString(base16, 16)
	if !ok {
		panic("bad number: " + base16)
	}
	return i
}

func SetPublicKey(a, b string) (n *big.Int, e int64) {
	if a != "" && b != "" && len(a) > 0 && len(b) > 0 {
		n = FromBase16(a)
		e, _ = strconv.ParseInt(b, 16, 0)
		return
	} else {
		fmt.Println("Invalid RSA public key")
	}

	return
}

func DoPublic(a int64) {

	// return a.modPowInt(this.e, this.n)

}
func encrypt(a string, n *big.Int, e int64) string {
	/*	var b = _(a, n.BitLen() + 7 >> 3);
		if (null == b)
		    return null;
		var c = this.doPublic(b);
		if (null == c)
		    return null;
		var d = c.toString(16);
		return 0 == (1 & d.length) ? d : "0" + d
	*/
	return ""

}

// 加密
func RsaEncrypt(publicKey []byte, origData []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}
