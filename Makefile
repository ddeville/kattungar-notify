##### Admin CLI #####

.PHONY: build-admin-cli run-admin-cli install-admin-cli

build-admin-cli:
	CGO_ENABLED=0 go build -o build/kattungar-notify-admin ./cmd/admin_cli

run-admin-cli:
	go run ./cmd/admin_cli

install-admin-cli: build-admin-cli
	sudo mv build/kattungar-notify-admin /usr/local/bin/kattungar-notify-admin

##### Notify CLI #####

.PHONY: build-notify-cli run-notify-cli install-notify-cli

build-notify-cli:
	CGO_ENABLED=0 go build -o build/kattungar-notify ./cmd/notify_cli

run-notify-cli:
	go run ./cmd/notify_cli

install-notify-cli: build-notify-cli
	sudo mv build/kattungar-notify /usr/local/bin/kattungar-notify

##### Server #####

.PHONY: build-server run-server

build-server:
	docker-compose -f docker-compose.yaml build

run-server:
	docker-compose -f docker-compose.yaml up

##### iOS #####

.PHONY: build-ios archive-ios publish-ios

build-ios:
	xcodebuild -project "ios/KattungarNotify.xcodeproj" -configuration Debug -scheme "Kattungar Notify" -allowProvisioningUpdates

archive-ios:
	security unlock-keychain
	rm -rf build/ios
	xcodebuild clean -project "ios/KattungarNotify.xcodeproj" -scheme "Kattungar Notify" -configuration Release -destination generic/platform=iOS -sdk iphoneos
	xcodebuild archive -project "ios/KattungarNotify.xcodeproj" -scheme "Kattungar Notify" -configuration Release -destination generic/platform=iOS -archivePath build/ios/KattungarNotify.xcarchive -allowProvisioningUpdates

publish-ios: archive-ios
	xcodebuild -exportArchive -archivePath build/ios/KattungarNotify.xcarchive -exportOptionsPlist ios/export_options_ios.plist -exportPath build/ios -allowProvisioningUpdates

##### macOS #####

build-macos:
	xcodebuild -project "ios/KattungarNotify.xcodeproj" -configuration Debug -scheme "Kattungar Notify MacOS" -allowProvisioningUpdates

archive-macos:
	security unlock-keychain
	rm -rf build/macos
	xcodebuild clean -project "ios/KattungarNotify.xcodeproj" -scheme "Kattungar Notify MacOS" -configuration Release -destination generic/platform=macOS -sdk macosx
	xcodebuild archive -project "ios/KattungarNotify.xcodeproj" -scheme "Kattungar Notify MacOS" -configuration Release -destination generic/platform=macOS -archivePath build/macos/KattungarNotify.xcarchive -allowProvisioningUpdates

publish-macos: archive-macos
	./scripts/export-and-notarize-macos
