# Webauthn Demo

W3C 规范: [webappsec-credential-management](https://w3c.github.io/webappsec-credential-management), [webauthn](https://www.w3.org/TR/webauthn/)

## 过程

- 必须用域名，必须有证书(自签名亦可)
- navigator.credentials.create() 参数对象结构必须正确，
  publicKey.user.id, publicKey.challenge 类型必须为 ArrayBuffer
- await navigator.credentials.create()
- Reset Keys: "Reset your security key" of chrome://settings/securityKeys

## gin 不如 fasthttp 好用

- gin.Context 设计太诡异
- gin.Context.{Get|Bind}
- middleware 机制较为复杂

## 术语

- RP: Relying Party
- CTAP: Client-to-Authenticator Protocols
- FQDN: Fully qualified domain name

## TODO

<!-- https://fidoalliance.org/specs/fido-v2.0-ps-20190130/fido-client-to-authenticator-protocol-v2.0-ps-20190130.html#ble-gatt-service -->

### Android

- 可能要 [ChromeCustomTabs](https://developer.chrome.com/docs/android/custom-tabs/) 才能支持。flutter 支持为 [flutter_custom_tabs](https://pub.dev/packages/flutter_custom_tabs).

- 使用 [Fido2ApiClient](https://developers.google.com/android/reference/com/google/android/gms/fido/fido2/Fido2ApiClient) [文档](https://developers.google.com/identity/fido/android/native-apps) [stackoverflow](https://stackoverflow.com/questions/57674215/how-to-implement-webauthn-in-an-android-app)
