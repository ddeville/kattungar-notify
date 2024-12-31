//
//  ApplicationDelegate.swift
//  Kattungar Notify
//
//  Created by Damien Deville on 1/2/24.
//

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
import UserNotifications

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

    @objc func application(_ application: Application, didRegisterForRemoteNotificationsWithDeviceToken deviceToken: Data) {
        let token = deviceToken.map { data in String(format: "%02.2hhx", data) }.joined()
        let deviceKey = UserDefaults.standard.string(forKey: DeviceKeyDefaultsKey)!

        // Remove existing token so that we attempt to refresh next time if we fail retrieving
        UserDefaults.standard.removeObject(forKey: TokenDefaultsKey)

        registerToken(deviceKey: deviceKey, token: token) { result in
            DispatchQueue.main.async {
                switch result {
                case .success(_):
                    print("Successfully updated token")
                    UserDefaults.standard.set(token, forKey:TokenDefaultsKey)
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
        print("Failed to register for notifications... \(error)")
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
