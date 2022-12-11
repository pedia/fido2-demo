# 目标
```text
          安卓(chrome) 指纹/Face
        /
Desktop -- iOS 指纹/Face
        \  
          本机 yubikey

iOS App -- 本机指纹/Face
```

| |认证方式|可实现|认证方式|体验|
|-----|-----|-----|-----|-----|
| Web | Safri 直接 WebAuthn 认证 |已实现|AppleID认证||
| Web | Chrome 直接 WebAuthn 认证 |已实现|多种方式不可控 key/本机指纹/手机Chrome||
| Desktop App | 内嵌 yubi SDK 直接 yubikey 认证 |可| key ||
| Desktop App | 嵌入 Chrome WebAuthn 认证|已实现|多种方式不可控 key/本机指纹/手机Chrome||
| Mac App | 打开 Safari |已实现|AppleID认证||
| iOS App | 打开 Safari |已实现|AppleID认证||
| iOS App | 内嵌 yubi SDK |可|key(nfc)|不好|
| Android App | 内嵌 yubi SDK 直接 yubikey 认证 |可 | key ||
| Android App | 内嵌 ChromeTab |可|key/本机指纹||


# Fido2 Demo

- U2F Chrome/Mozilla 实现的 Webauthn 规范
- OTP yubicloud 实现的一种 passwordless
- OATH is an organization that specifies two open authentication standards: TOTP and HOTP.

Webauthn(W3C 规范): [webappsec-credential-management](https://w3c.github.io/webappsec-credential-management), [webauthn](https://www.w3.org/TR/webauthn/)

## Webauthn 实现过程

- 必须用域名，必须有证书(自签名亦可)
- navigator.credentials.create() 参数对象结构必须正确，
  publicKey.user.id, publicKey.challenge 类型必须为 ArrayBuffer
- await navigator.credentials.create()
- Reset Keys: "Reset your security key" of chrome://settings/securityKeys
- 错误使用 base64URLStringToBuffer/stringToBuffer 浪费半天才解决

### gin 不如 fasthttp 好用

- gin.Context 设计太诡异
- gin.Context.{Get|Bind}
- middleware 机制较为复杂

### 术语

- RP: Relying Party
- CTAP: Client-to-Authenticator Protocols
- FQDN: Fully qualified domain name

### TODO

#### Android

- 可能要 [ChromeCustomTabs](https://developer.chrome.com/docs/android/custom-tabs/) 才能支持。flutter 支持为 [flutter_custom_tabs](https://pub.dev/packages/flutter_custom_tabs).

- 使用 [Fido2ApiClient](https://developers.google.com/android/reference/com/google/android/gms/fido/fido2/Fido2ApiClient) [文档](https://developers.google.com/identity/fido/android/native-apps) [stackoverflow](https://stackoverflow.com/questions/57674215/how-to-implement-webauthn-in-an-android-app)
