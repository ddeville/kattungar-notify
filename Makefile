.PHONY: build-cli run-cli install-cli build-server run-server build-ios archive-ios publish-ios

build-cli:
	go build -o build/kattungar-notify-admin ./cmd/cli

run-cli:
	go run ./cmd/cli

install-cli: build-cli
	sudo mv build/kattungar-notify-admin /usr/local/bin/kattungar-notify-admin

build-server:
	docker-compose -f docker-compose.yaml build

run-server:
	docker-compose -f docker-compose.yaml up

build-ios:
	xcodebuild -project "ios/KattungarNotify.xcodeproj" -configuration Debug -scheme "Kattungar Notify" -allowProvisioningUpdates

archive-ios:
	rm -rf build/ios
	xcodebuild clean -project "ios/KattungarNotify.xcodeproj" -scheme "Kattungar Notify" -configuration Release -destination generic/platform=iOS -sdk iphoneos
	xcodebuild archive -project "ios/KattungarNotify.xcodeproj" -scheme "Kattungar Notify" -configuration Release -destination generic/platform=iOS -archivePath build/ios/KattungarNotify.xcarchive -allowProvisioningUpdates

publish-ios: archive-ios
	xcodebuild -exportArchive -archivePath build/ios/KattungarNotify.xcarchive -exportOptionsPlist ios/exportOptions.plist -exportPath build/ios -allowProvisioningUpdates
