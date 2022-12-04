# Webauthn Demo

## 过程

- navigator.credentials.create() 参数对象结构必须正确，
  publicKey.user.id, publicKey.challenge 类型必须为 ArrayBuffer
- 必须用域名，必须有证书(自签名亦可)
- await navigator.credentials.create()

## gin 不如 fasthttp 好用

- gin.Context 设计太诡异
- gin.Context.Get
- gin.Context.Bind
- middleware 机制较为复杂

## 术语

- RP: Relying Party
- CTAP: Client-to-Authenticator Protocols
- FQDN: Fully qualified domain name
