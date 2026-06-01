//
//  ApplicationDelegate.swift
//  Kattungar Notify
//
//  Created by Damien Deville on 1/2/24.
//

import UserNotifications

#if os(iOS)

import UIKit

typealias Application = UIApplication

class ApplicationDelegate: CommonApplicationDelegate, UIApplicationDelegate {
    func application(_ application: UIApplication, didFinishLaunchingWithOptions launchOptions: [UIApplication.LaunchOptionsKey: Any]? = nil) -> Bool {
        registerForRemoteNotifications = application.registerForRemoteNotifications
        handleApplicationChange()
        return true
    }
}

#elseif os(macOS)

import Cocoa

typealias Application = NSApplication

class ApplicationDelegate: CommonApplicationDelegate, NSApplicationDelegate {
    func applicationDidFinishLaunching(_ notification: Notification) {
        registerForRemoteNotifications = (notification.object as! NSApplication).registerForRemoteNotifications
        handleApplicationChange()
    }
}

#endif

class CommonApplicationDelegate: NSObject, ObservableObject {
    fileprivate var registerForRemoteNotifications: (() -> Void)?

    @Published var hasSetupDeviceKey: Bool = false
    @Published var isSendingToken: Bool = false

    private var shouldNotifyAfterManualTokenSend = false

    func handleApplicationChange() {
        UNUserNotificationCenter.current().delegate = self

        // This will switch between the device key setup and regular view
        hasSetupDeviceKey = UserDefaults.standard.string(forKey: DeviceKeyDefaultsKey) != nil

        // Let's get a token if we've never retrieved one before
        if hasSetupDeviceKey && UserDefaults.standard.string(forKey: TokenDefaultsKey) == nil {
            registerForRemoteNotifications!()

            UNUserNotificationCenter.current().requestAuthorization(options: [.alert, .badge, .sound]) { success, error in
                if !success {
                    print("Didn't get approval to send push notifications... \(String(describing: error))")
                }
            }
        }
    }

    func manuallySendTokenToServer() {
        guard !isSendingToken else {
            return
        }
        guard let deviceKey = UserDefaults.standard.string(forKey: DeviceKeyDefaultsKey) else {
            hasSetupDeviceKey = false
            return
        }

        isSendingToken = true

        if let token = UserDefaults.standard.string(forKey: TokenDefaultsKey) {
            sendTokenToServer(deviceKey: deviceKey, token: token, isManualTokenSend: true)
            return
        }

        guard let registerForRemoteNotifications else {
            isSendingToken = false
            print("Cannot register for remote notifications before the application has finished launching")
            return
        }

        shouldNotifyAfterManualTokenSend = true
        registerForRemoteNotifications()
    }

    @objc func application(_ application: Application, didRegisterForRemoteNotificationsWithDeviceToken deviceToken: Data) {
        let token = deviceToken.map { data in String(format: "%02.2hhx", data) }.joined()
        let deviceKey = UserDefaults.standard.string(forKey: DeviceKeyDefaultsKey)!
        let isManualTokenSend = shouldNotifyAfterManualTokenSend

        // Remove existing token so that we attempt to refresh next time if we fail retrieving
        UserDefaults.standard.removeObject(forKey: TokenDefaultsKey)

        sendTokenToServer(deviceKey: deviceKey, token: token, isManualTokenSend: isManualTokenSend)
    }

    private func sendTokenToServer(deviceKey: String, token: String, isManualTokenSend: Bool) {
        registerToken(deviceKey: deviceKey, token: token) { result in
            DispatchQueue.main.async {
                if isManualTokenSend {
                    self.isSendingToken = false
                    self.shouldNotifyAfterManualTokenSend = false
                }

                switch result {
                case .success(let (response, _)):
                    print("Successfully updated token")
                    UserDefaults.standard.set(token, forKey:TokenDefaultsKey)
                    if isManualTokenSend && response.statusCode == 200 {
                        self.showTokenSentNotification()
                    }
                case .failure(let error):
                    if case URLSession.HTTPError.serverSideError(let statusCode) = error {
                        if statusCode == 401 {
                            print("Request failed as unauthorized, device key is likely wrong: \(error)")
                            UserDefaults.standard.removeObject(forKey: TokenDefaultsKey)
                            self.hasSetupDeviceKey = false
                        }
                    }
                    print("Failed to make request to server \(error)")
                }
            }
        }
    }

    @objc func application(_ application: Application, didFailToRegisterForRemoteNotificationsWithError error: Error) {
        if shouldNotifyAfterManualTokenSend {
            isSendingToken = false
            shouldNotifyAfterManualTokenSend = false
        }
        print("Failed to register for notifications... \(error)")
    }

    private func showTokenSentNotification() {
        let content = UNMutableNotificationContent()
        content.title = "Kattungar Notify"
        content.body = "Token sent to server."
        content.sound = .default

        let request = UNNotificationRequest(identifier: "token-sent-\(UUID().uuidString)", content: content, trigger: nil)
        UNUserNotificationCenter.current().add(request) { error in
            if let error = error {
                print("Failed to show token sent notification: \(error)")
            }
        }
    }
}

extension CommonApplicationDelegate: UNUserNotificationCenterDelegate {
    func userNotificationCenter(_ center: UNUserNotificationCenter, willPresent notification: UNNotification, withCompletionHandler completionHandler: @escaping (UNNotificationPresentationOptions) -> Void) {
        completionHandler([[.banner, .sound]])
    }

    func userNotificationCenter(_ center: UNUserNotificationCenter, didReceive response: UNNotificationResponse, withCompletionHandler completionHandler: @escaping () -> Void) {
        completionHandler()
    }
}
