package dsl

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"gitlab.shudieds.com/mec/lib/utils/convert"
	"gitlab.shudieds.com/zxh/engine/types"
	"golang.org/x/crypto/bcrypt"
	"io"
)

/*
*
base64Encode
*/
func (a *AST) strBase64Encode(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	str, bl := convert.String(v1)
	if bl == false {
		return nil, errors.New("strBase64Encode 参数错误")
	}
	return base64.StdEncoding.EncodeToString([]byte(str)), nil
}

/*
*

	base64Decode
*/
func (a *AST) strBase64Decode(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	str, bl := convert.String(v1)
	if bl == false {
		return nil, errors.New("strBase64Decode 参数错误")
	}

	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, errors.New("strBase64Decode err:" + err.Error())
	}
	return string(data), nil
}

/*
*

	md5
*/
func (a *AST) strMd5(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	str, bl := convert.String(v1)
	if bl == false {
		return nil, errors.New("strMd5 参数错误")
	}

	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil)), nil
}

/*
*

	hmac
*/
func (a *AST) strHmac(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}
	str, bl1 := convert.String(v1)
	key, bl2 := convert.String(v2)
	if bl1 == false || bl2 == false {
		return nil, errors.New("strHmac 参数错误")
	}

	hash := hmac.New(md5.New, []byte(key)) // 创建对应的md5哈希加密算法
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum([]byte(""))), nil
}

/*
*

	HmacSha256
*/
func (a *AST) strHmacSha256(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}
	str, bl1 := convert.String(v1)
	key, bl2 := convert.String(v2)
	if bl1 == false || bl2 == false {
		return nil, errors.New("strHmacSha256 参数错误")
	}

	hash := hmac.New(sha256.New, []byte(key)) // 创建对应的md5哈希加密算法
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum([]byte(""))), nil
}

/*
*

	sha1
*/
func (a *AST) strSha1(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	str, bl := convert.String(v1)
	if bl == false {
		return nil, errors.New("strSha1 参数错误！")
	}
	sha := sha1.New()
	sha.Write([]byte(str))
	return hex.EncodeToString(sha.Sum([]byte(""))), nil
}

/*
*

	sha256
*/
func (a *AST) strSha256(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	str, bl := convert.String(v1)
	if bl == false {
		return nil, errors.New("strSha256 参数错误！")
	}
	sha := sha256.New()
	sha.Write([]byte(str))
	return hex.EncodeToString(sha.Sum([]byte(""))), nil
}

/*
*
加盐加密
*/
func (a *AST) strEncrypt(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	str, bl := convert.String(v1)
	if bl == false {
		return nil, errors.New("strEncrypt 参数错误！")
	}

	hashStr, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("strEncrypt err:" + err.Error())
	}
	return string(hashStr), nil
}

/*
*
aes ECB 加密
*/
func (a *AST) strAesECBEncode(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}

	source, bl1 := convert.String(v1)
	key, bl2 := convert.String(v2)
	if bl1 == false || bl2 == false {
		return nil, errors.New("strAesECBEncode 参数错误!")
	}

	encryptCode := aesEncryptECB([]byte(source), []byte(key))
	return string(encryptCode), nil
}

/*
*
aes ECB 解密
*/
func (a *AST) strAesECBDecode(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}

	source, bl1 := convert.String(v1)
	key, bl2 := convert.String(v2)
	if bl1 == false || bl2 == false {
		return nil, errors.New("strAesECBDecode 参数错误!")
	}
	decryptCode := aesDecryptECB([]byte(source), []byte(key))
	return string(decryptCode), nil

}

/*
*
aes CBC 加密
*/
func (a *AST) strAesCBCEncode(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}

	source, bl1 := convert.String(v1)
	key, bl2 := convert.String(v2)
	if bl1 == false || bl2 == false {
		return nil, errors.New("strAesCBCEncode 参数错误!")
	}
	encryptCode := aesCBCEncrypt(source, key)
	return encryptCode, nil
}

/*
*
aes CBC 解密
*/
func (a *AST) strAesCBCDecode(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}

	source, bl1 := convert.String(v1)
	key, bl2 := convert.String(v2)
	if bl1 == false || bl2 == false {
		return nil, errors.New("strAesCBCDecode 参数错误!")
	}

	decryptCode := AesCBCDecrypt(source, key)
	return decryptCode, nil

}

/*
*
aes CRT 加密
*/
func (a *AST) strAesCRTEncode(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}

	source, bl1 := convert.String(v1)
	key, bl2 := convert.String(v2)
	if bl1 == false || bl2 == false {
		return nil, errors.New("strAesCRTEncode 参数错误!")
	}

	encryptCode, err := aesCtrCrypt([]byte(source), []byte(key))
	if err != nil {
		return nil, errors.New("strAesCRTEncode err:" + err.Error())
	}
	return string(encryptCode), nil
}

/*
*
aes CRT 解密
*/
func (a *AST) strAesCRTDecode(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}

	source, bl1 := convert.String(v1)
	key, bl2 := convert.String(v2)
	if bl1 == false || bl2 == false {
		return nil, errors.New("strAesCRTDecode 参数错误!")
	}

	encryptCode, err := aesCtrCrypt([]byte(source), []byte(key))
	if err != nil {
		return nil, errors.New("strAesCRTDecode err:" + err.Error())
	}
	return string(encryptCode), nil
}

/*
*
aes CFB 加密
*/
func (a *AST) strAesCFBEncode(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}

	source, bl1 := convert.String(v1)
	key, bl2 := convert.String(v2)
	if bl1 == false || bl2 == false {
		return nil, errors.New("strAesCFBEncode 参数错误!")
	}

	encryptCode := aesEncryptCFB([]byte(source), []byte(key))
	return string(encryptCode), nil
}

/*
*
aes CFB 解密
*/
func (a *AST) strAesCFBDecode(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}

	source, bl1 := convert.String(v1)
	key, bl2 := convert.String(v2)
	if bl1 == false || bl2 == false {
		return nil, errors.New("strAesCFBDecode 参数错误!")
	}

	encryptCode := aesDecryptCFB([]byte(source), []byte(key))
	return string(encryptCode), nil
}

/*
*
aes OFB 加密
*/
func (a *AST) strAesOFBEncode(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}

	source, bl1 := convert.String(v1)
	key, bl2 := convert.String(v2)
	if bl1 == false || bl2 == false {
		return nil, errors.New("strAesOFBEncode 参数错误!")
	}

	encryptCode, err := aesEncryptOFB([]byte(source), []byte(key))
	if err != nil {
		return nil, errors.New("strAesOFBEncode err:" + err.Error())
	}
	return string(encryptCode), nil
}

/*
*
aes OFB 解密
*/
func (a *AST) strAesOFBDecode(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}

	source, bl1 := convert.String(v1)
	key, bl2 := convert.String(v2)
	if bl1 == false || bl2 == false {
		return nil, errors.New("strAesOFBDecode 参数错误!")
	}
	encryptCode, err := aesDecryptOFB([]byte(source), []byte(key))
	if err != nil {
		return nil, errors.New("strAesOFBDecode err:" + err.Error())
	}
	return string(encryptCode), nil
}

// =================== ECB ======================
func aesEncryptECB(origData []byte, key []byte) (encrypted []byte) {
	cipher, _ := aes.NewCipher(generateKey(key))
	length := (len(origData) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, origData)
	pad := byte(len(plain) - len(origData))
	for i := len(origData); i < len(plain); i++ {
		plain[i] = pad
	}
	encrypted = make([]byte, len(plain))
	// 分组分块加密
	for bs, be := 0, cipher.BlockSize(); bs <= len(origData); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Encrypt(encrypted[bs:be], plain[bs:be])
	}

	return encrypted
}
func aesDecryptECB(encrypted []byte, key []byte) (decrypted []byte) {
	cipher, _ := aes.NewCipher(generateKey(key))
	decrypted = make([]byte, len(encrypted))
	//
	for bs, be := 0, cipher.BlockSize(); bs < len(encrypted); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}

	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}

	return decrypted[:trim]
}
func generateKey(key []byte) (genKey []byte) {
	genKey = make([]byte, 16)
	copy(genKey, key)
	for i := 16; i < len(key); {
		for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	return genKey
}

/** CBC **/

func aesCBCEncrypt(orig string, key string) string {
	// 转成字节数组
	origData := []byte(orig)
	k := []byte(key)
	// 分组秘钥
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	block, _ := aes.NewCipher(k)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 补全码
	origData = pkcS7Padding(origData, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	// 创建数组
	cryted := make([]byte, len(origData))
	// 加密
	blockMode.CryptBlocks(cryted, origData)
	return base64.StdEncoding.EncodeToString(cryted)
}
func AesCBCDecrypt(cryted string, key string) string {
	// 转成字节数组
	crytedByte, _ := base64.StdEncoding.DecodeString(cryted)
	k := []byte(key)
	// 分组秘钥
	block, _ := aes.NewCipher(k)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	// 创建数组
	orig := make([]byte, len(crytedByte))
	// 解密
	blockMode.CryptBlocks(orig, crytedByte)
	// 去补全码
	orig = pkcS7UnPadding(orig)
	return string(orig)
}

// 补码
// AES加密数据块分组长度必须为128bit(byte[16])，密钥长度可以是128bit(byte[16])、192bit(byte[24])、256bit(byte[32])中的任意一个。
func pkcS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

/**CRT模式**/
//加密
func aesCtrCrypt(plainText []byte, key []byte) ([]byte, error) {

	//1. 创建cipher.Block接口
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//2. 创建分组模式，在crypto/cipher包中
	iv := bytes.Repeat([]byte("1"), block.BlockSize())
	stream := cipher.NewCTR(block, iv)
	//3. 加密
	dst := make([]byte, len(plainText))
	stream.XORKeyStream(dst, plainText)

	return dst, nil
}

/** CFB **/
func aesEncryptCFB(origData []byte, key []byte) (encrypted []byte) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil
	}
	encrypted = make([]byte, aes.BlockSize+len(origData))
	iv := encrypted[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(encrypted[aes.BlockSize:], origData)
	return encrypted
}
func aesDecryptCFB(encrypted []byte, key []byte) (decrypted []byte) {
	block, _ := aes.NewCipher(key)
	if len(encrypted) < aes.BlockSize {
		return nil
	}
	iv := encrypted[:aes.BlockSize]
	encrypted = encrypted[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(encrypted, encrypted)
	return encrypted
}

/**OFB**/
func aesEncryptOFB(data []byte, key []byte) ([]byte, error) {
	data = pkcS7Padding(data, aes.BlockSize)
	block, _ := aes.NewCipher([]byte(key))
	out := make([]byte, aes.BlockSize+len(data))
	iv := out[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewOFB(block, iv)
	stream.XORKeyStream(out[aes.BlockSize:], data)
	return out, nil
}

func aesDecryptOFB(data []byte, key []byte) ([]byte, error) {
	block, _ := aes.NewCipher([]byte(key))
	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]
	if len(data)%aes.BlockSize != 0 {
		return nil, errors.New("data is not a multiple of the block size")
	}

	out := make([]byte, len(data))
	mode := cipher.NewOFB(block, iv)
	mode.XORKeyStream(out, data)

	out = pkcS7UnPadding(out)
	return out, nil
}
