package chat

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"
)

func pad(src []byte) []byte {
	padding := aes.BlockSize - len(src)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func unpad(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])

	if unpadding > length {
		panic("unpad error. This could happen when incorrect encryption key is used")
	}

	return src[:(length - unpadding)]
}

func theAESencrypter(skey []byte, data []byte) []byte {
	plaintext := []byte(pad(data))
	block, err := aes.NewCipher(skey[:16])
	if err != nil {
		panic(err)
	}
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)
	return ciphertext
}

func theAESdecrypter(skey []byte, ciphertext []byte) []byte {
	block, err := aes.NewCipher(skey[:16])
	if err != nil {
		panic(err)
	}
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ct := ciphertext[aes.BlockSize:]

	if len(ct)%aes.BlockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	result := make([]byte, len(ct))
	mode.CryptBlocks(result, ct)
	return result
}

func base64Encoder(ciphertext []byte) string {
	s := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(ciphertext)
	return s
}

func base64Decoder(ciphertext string) []byte {
	bs, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(ciphertext)
	if err != nil {
		panic("Decode error")
		error.Error(err)
	}
	return bs
}

func Encrypt(app_secret string, data *map[string]interface{}) string {
	jdata, _ := json.Marshal(data)
	skey := sha256.Sum256([]byte(app_secret))
	r := theAESencrypter(skey[:], jdata)
	s := base64Encoder(r)
	return s
}

func Decrypt(app_secret string, data_str string, data_ptr *map[string]interface{}) {
	bs := base64Decoder(data_str)
	skey := sha256.Sum256([]byte(app_secret))
	r := theAESdecrypter(skey[:], bs)
	unpadMsg := unpad(r)
	err := json.Unmarshal(unpadMsg, data_ptr)
	if err != nil {
		// panic("JSON error")
		error.Error(err)
	}
}

func Decrypt2String(app_secret string, data_str string) string {
	bs := base64Decoder(data_str)
	skey := sha256.Sum256([]byte(app_secret))
	r := theAESdecrypter(skey[:], bs)
	unpadMsg := unpad(r)
	log.Println("unpadMsg: ", unpadMsg)
	rst := string(unpadMsg[:])
	return rst
}

func testAes() {
	app_secret := "123321"
	data := map[string]interface{}{
		"timestamp": time.Now().UTC().Unix(),
		"user_id":   1,
		"img":       "http://xxx.com",
	}
	encrypted := Encrypt(app_secret, &data)
	data2 := map[string]interface{}{}
	Decrypt(app_secret, encrypted, &data2)
	if data["img"] != data2["img"] {
		panic("aes failed")
	}
	fmt.Println(app_secret, encrypted)
}
