client.dmg: build/macos/Build/Products/Release/client.app
	hdiutil create -volname Fido2Demo \
		-srcfolder build/macos/Build/Products/Release/client.app \
		-ov -format UDZO fido2demo.dmg

# create-dmg \
# 	--volname "Yubico Authenticator" \
# 	--background "dmg-background.png" \
# 	--window-pos 200 120 \
# 	--window-size 800 400 \
# 	--icon-size 100 \
# 	--icon "Yubico Authenticator.app" 200 190 \
# 	--app-drop-link 600 185 \
# 	"yubioath-desktop.dmg" \