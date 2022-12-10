// {{ define "js" }}
// convert Unicode string to a string in which
// each 16-bit unit occupies only one byte
function toBinary(string) {
    const codeUnits = Uint16Array.from(
        { length: string.length },
        (element, index) => string.charCodeAt(index)
    );
    const charCodes = new Uint8Array(codeUnits.buffer);

    let result = "";
    charCodes.forEach((char) => {
        result += String.fromCharCode(char);
    });
    return result;
}

function base64URLStringToBuffer(base64URLString) {
    // atok('6ZxDKgiBAWc87KJLqZ38NQ==', 'base64')
    return Uint8Array.from(window.atob(base64URLString), c => c.charCodeAt(0));
}

function bufferToBase64URLString(buffer) {
    const bytes = new Uint8Array(buffer);
    let str = '';

    for (const charCode of bytes) {
        str += String.fromCharCode(charCode);
    }
    const base64String = btoa(str);
    return base64String.replace(/\+/g, '-').replace(/\//g, '_').replace(/=/g, '');
}

async function is_pkc_available() {
    return await PublicKeyCredential.isUserVerifyingPlatformAuthenticatorAvailable();
}

//
async function register(opts) {
    if (!window.PublicKeyCredential) {
        // Client not capable. Handle error.
        pageLog('warn', 'navigator.credentials not supported');
        return;
    }

    pageLog('navigator.credentials.create', opts)

    // Converto Json str to ArrayBuffer
    opts.publicKey.challenge = base64URLStringToBuffer(opts.publicKey.challenge);
    opts.publicKey.user.id = base64URLStringToBuffer(opts.publicKey.user.id);

    navigator.credentials.create(opts)
        .then(handle_register)
        .catch(handle_failed);
}

// Send new credential info to server for verification and registration.
function handle_register(credential) {
    // TODO: pageLog not worked for `credential`
    console.log('got raw', credential);

    const { id, rawId, response, type } = credential;
    console.log(id, rawId, response, type);

    const c = {
        id,
        rawId: bufferToBase64URLString(rawId),
        response: {
            attestationObject: bufferToBase64URLString(response.attestationObject),
            clientDataJSON: bufferToBase64URLString(response.clientDataJSON),
        },
        type,
        clientExtensionResults: credential.getClientExtensionResults(),
        authenticatorAttachment: credential.authenticatorAttachment,
    };
    send('/register', c);
}

// No acceptable authenticator or user refused consent. Handle appropriately.
function handle_failed(err) {
    console.log('failed', err);
    pageLog('failed', err);
}

async function login(opts) {
    if (!window.PublicKeyCredential) {
        pageLog('warn', 'navigator.credentials not supported');
        return;
    }

    pageLog('navigator.credentials.get', opts)

    // Converto Json str to ArrayBuffer
    opts.publicKey.challenge = base64URLStringToBuffer(opts.publicKey.challenge);
    for (var i = 0; i < opts.publicKey.allowCredentials.length; i++) {
        opts.publicKey.allowCredentials[i].id
            = base64URLStringToBuffer(opts.publicKey.allowCredentials[i].id);
    }

    navigator.credentials.get(opts)
        .then(handle_login)
        .catch(handle_failed);
}

function handle_login(credential) {
    // TODO: pageLog not worked for `credential`
    console.log('got raw', credential)

    const { id, rawId, response, type } = credential;
    // console.log(id, rawId, response, type);

    const c = {
        id,
        rawId: bufferToBase64URLString(rawId),
        response: {
            authenticatorData: bufferToBase64URLString(response.authenticatorData),
            clientDataJSON: bufferToBase64URLString(response.clientDataJSON),
            signature: bufferToBase64URLString(response.signature),
            userHandle: bufferToBase64URLString(response.userHandle),
        },
        type,
        clientExtensionResults: credential.getClientExtensionResults(),
        authenticatorAttachment: credential.authenticatorAttachment,
    };
    send('/login', c);
}

function userNameChanged() {
    var un = document.querySelector('#un').value;

    document.querySelector('#register').href = '/register?username=' + un;
    document.querySelector('#login').href = '/login?username=' + un;
}

function send(url, credential) {
    var xhr = new XMLHttpRequest();
    xhr.open("POST", url, true);
    xhr.setRequestHeader("Content-Type", "application/json; charset=utf-8");

    xhr.onreadystatechange = () => {
        if (xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200) { 
            pageLog('server response ' + xhr.status, xhr.responseText);
        }
    };

    var o = JSON.stringify(credential);
    xhr.send(o);
    pageLog('send to server ' + url, o);
}

function pageLog(ctx, message) {
    var s = '<b>' + ctx + '</b><br />\n';

    if (typeof message == 'object') {
        s +=  JSON.stringify(message) + '<br />\n';
    } else {
        s += message + '<br />\n';
    }

    var logger = document.getElementById('log');
    logger.innerHTML += s;
}
// {{ end }}
