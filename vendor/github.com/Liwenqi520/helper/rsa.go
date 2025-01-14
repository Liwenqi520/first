package helper

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

//rsa writer
type RsaStore struct {
	Data []byte
}

func (r *RsaStore) Write(p []byte) (n int, err error) {
	r.Data = append(r.Data, p...)
	return
}

type RSA struct {
}

//生成RSA私钥和公钥，保存到文件中
func (c RSA) GenerateRSAKey(bits int) (public, private []byte, err error) {
	//GenerateKey函数使用随机数据生成器random生成一对具有指定字位数的RSA密钥
	//Reader是一个全局、共享的密码用强随机数生成器
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return public, private, err
	}
	//保存私钥
	//通过x509标准将得到的ras私钥序列化为ASN.1 的 DER编码字符串
	X509PrivateKey := x509.MarshalPKCS1PrivateKey(privateKey)
	//使用pem格式对x509输出的内容进行编码
	//构建一个pem.Block结构体对象
	privateBlock := pem.Block{Type: "PRIVATE KEY", Bytes: X509PrivateKey}
	//将数据保存
	rsaStorePrivate := RsaStore{}
	err = pem.Encode(&rsaStorePrivate, &privateBlock)
	if err != nil {
		return public, private, err
	}
	//保存公钥
	//获取公钥的数据
	publicKey := privateKey.PublicKey
	//X509对公钥编码
	X509PublicKey, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		return public, private, err
	}
	//pem格式编码
	//创建一个pem.Block结构体对象
	publicBlock := pem.Block{Type: "PUBLIC KEY", Bytes: X509PublicKey}
	//保存
	rsaStorePublic := &RsaStore{}
	err = pem.Encode(rsaStorePublic, &publicBlock)
	if err != nil {
		return public, private, err
	}
	return rsaStorePublic.Data, rsaStorePrivate.Data, nil
}

//RSA加密
func (c RSA) RSAEncrypt(plainText []byte, publicBuf []byte) (b []byte, err error) {
	//pem解码
	block, _ := pem.Decode(publicBuf)
	//x509解码
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return b, err
	}
	//类型断言
	publicKey := publicKeyInterface.(*rsa.PublicKey)
	//对明文进行加密
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plainText)
	if err != nil {
		return b, err
	}
	//返回密文
	return cipherText, nil
}

//RSA解密
func (c RSA) RSADecrypt(cipherText []byte, privateBuf []byte) (b []byte, err error) {
	//pem解码
	block, _ := pem.Decode(privateBuf)
	//X509解码
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return b, err
	}
	//对密文进行解密
	plainText, _ := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherText)
	//返回明文
	return plainText, nil
}
