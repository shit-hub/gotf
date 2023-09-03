package encrypt_test

import (
	"encoding/hex"
	"fmt"

	"gotf/encrypt"
)

func ExampleAesEncryptCBC() {
	origData := []byte("Hello World")
	key := []byte("ABCDEFGHIJKLMNOP")

	encrypted := encrypt.AesEncryptCBC(origData, key)
	fmt.Println("encrypt(hex): ", hex.EncodeToString(encrypted))
	decrypted := encrypt.AesDecryptCBC(encrypted, key)
	fmt.Println("decrypt: ", string(decrypted))

	//Output:
	//encrypt(hex):  20c7cbfe27e4baf919c06fefd9b9fa07
	//decrypt:  Hello World
}

func ExampleAesEncryptECB() {
	origData := []byte("Hello World")
	key := []byte("ABCDEFGHIJKLMNOP")

	encrypted := encrypt.AesEncryptECB(origData, key)
	fmt.Println("encrypt(hex): ", hex.EncodeToString(encrypted))
	decrypted := encrypt.AesDecryptECB(encrypted, key)
	fmt.Println("decrypt: ", string(decrypted))

	//Output:
	//encrypt(hex):  dc6e0e53ad9d33dcd78646d7cbca66b8
	//decrypt:  Hello World
}

func ExampleAesEncryptCFB() {
	origData := []byte("Hello World")
	key := []byte("ABCDEFGHIJKLMNOP")

	encrypted := encrypt.AesEncryptCFB(origData, key)
	fmt.Println("encrypt(hex): ", hex.EncodeToString(encrypted))
	decrypted := encrypt.AesDecryptCFB(encrypted, key)
	fmt.Println("decrypt: ", string(decrypted))
}
