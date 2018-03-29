package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	//	"io"
	"math/big"
	"testing"
)

var (
	rsaPrivateKey = &rsa.PrivateKey{
		PublicKey: rsa.PublicKey{
			N: FromBase16("9a39c3fefeadf3d194850ef3a1d707dfa7bec0609a60bfcc7fe4ce2c615908b9599c8911e800aff684f804413324dc6d9f982f437e95ad60327d221a00a2575324263477e4f6a15e3b56a315e0434266e092b2dd5a496d109cb15875256c73a2f0237c5332de28388693c643c8764f137e28e8220437f05b7659f58c4df94685"),
			E: parseInt16("10001"),
		},

		// 以下暂时没用
		D: fromBase10("7266398431328116344057699379749222532279343923819063639497049039389899328538543087657733766554155839834519529439851673014800261285757759040931985506583861"),
		Primes: []*big.Int{
			fromBase10("98920366548084643601728869055592650835572950932266967461790948584315647051443"),
			fromBase10("94560208308847015747498523884063394671606671904944666360068158221458669711639"),
		},
	}
)

func TestEncrypt(t *testing.T) {
	pwd := []byte("13640500hZm")
	random := rand.Reader

	//加密
	k := (rsaPrivateKey.N.BitLen() + 7) / 8
	if len(pwd) > k-11 {
		pwd = pwd[0 : k-11]
	}
	ciphertext, err := rsa.EncryptPKCS1v15(random, &rsaPrivateKey.PublicKey, pwd)
	if err != nil {
		fmt.Println("error encrypting: %s", err)
		//return false
	}

	fmt.Println("result", hex.EncodeToString(ciphertext), err)

	/*
		var publicKeyData = `-----BEGIN PUBLIC KEY-----
			d3bcef1f00424f3261c89323fa8cdfa12bbac400d9fe8bb627e8d27a44bd5d59dce559135d678a8143beb5b8d7056c4e1f89c4e1f152470625b7b41944a97f02da6f605a49a93ec6eb9cbaf2e7ac2b26a354ce69eb265953d2c29e395d6d8c1cdb688978551aa0f7521f290035fad381178da0bea8f9e6adce39020f513133fb
			-----END PUBLIC KEY-----
			`
		lM := "d3bcef1f00424f3261c89323fa8cdfa12bbac400d9fe8bb627e8d27a44bd5d59dce559135d678a8143beb5b8d7056c4e1f89c4e1f152470625b7b41944a97f02da6f605a49a93ec6eb9cbaf2e7ac2b26a354ce69eb265953d2c29e395d6d8c1cdb688978551aa0f7521f290035fad381178da0bea8f9e6adce39020f513133fb"
		lN, lE := SetPublicKey(lM, "10001")

		//lN = `-----BEGIN PUBLIC KEY-----` + lN + `-----END PUBLIC KEY-----`
		fmt.Println("SetPublicKey", FromBase16(lM), lN, lE)
		re, err := RsaEncrypt([]byte(publicKeyData), []byte("13640500hZm"))
		fmt.Println("RsaEncrypt", re, err)
		block, _ := pem.Decode([]byte(publicKeyData))
		fmt.Println("SetPublicKey", block)
		hash := sha1.New()
		random := rand.Reader
		msg := []byte("13640500hZm")
		var pub *rsa.PublicKey
		pubInterface, parseErr := x509.ParsePKIXPublicKey(block.Bytes)
		if parseErr != nil {
			fmt.Println("Load public key error")
			panic(parseErr)
		}
		pub = pubInterface.(*rsa.PublicKey)
		encryptedData, encryptErr := rsa.EncryptOAEP(hash, random, pub, msg, nil)
		if encryptErr != nil {
			fmt.Println("Encrypt data error")
			panic(encryptErr)
		}
		encodedData := base64.URLEncoding.EncodeToString(encryptedData)
		fmt.Println(encodedData)
	*/
}

func TestDisencrypt(t *testing.T) {
	//解密
	lpwd2, _ := hex.DecodeString("18c9570c06538c181f56331fbce882c9c9c5d89844b16ee91c9113c5ac5d311c8113da182e4a784633ecc30a09cbca9fb2fd7fb2d83fcbb6994cc7e5771d7e241a73ae521f31f54619563f4316177496ae349313eb7d933f732713d1da7be441909c9e5dbc8caf9e1105a6fa47fc558ab44089b7f9079ae64767b44ec0508f9a")
	random := rand.Reader
	fmt.Println("result", lpwd2)
	plaintext, err := rsa.DecryptPKCS1v15(random, rsaPrivateKey, lpwd2)
	if err != nil {
		fmt.Println("error decrypting: %s", err)
		//return false
	}
	fmt.Println("result", plaintext, err)
}
