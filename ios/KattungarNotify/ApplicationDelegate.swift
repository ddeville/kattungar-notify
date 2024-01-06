//
//  ApplicationDelegate.swift
//  Kattungar Notify
//
//  Created by Damien Deville on 1/2/24.
//

#if os(iOS)

import UIKit

typealias Application = UIApplication
typealias RemoteNotificationDelegate = UIApplicationDelegate

class ApplicationDelegate: UIResponder, UIApplicationDelegate, ObservableObject {
    func application(_ application: UIApplication, didFinishLaunchingWithOptions launchOptions: [UIApplication.LaunchOptionsKey: Any]? = nil) -> Bool {
        setupNotifications(application: application, delegate: self)
        return true
    }
}

#elseif os(macOS)

import Cocoa
import UserNotifications
import SwiftUI

typealias Application = NSApplication
typealias RemoteNotificationDelegate = NSApplicationDelegate

class ApplicationDelegate: NSObject, NSApplicationDelegate, ObservableObject {
    func applicationDidFinishLaunching(_ notification: Notification) {
        setupNotifications(application: notification.object as! NSApplication, delegate: self)
    }
}

#endif

func setupNotifications(application: Application, delegate: UNUserNotificationCenterDelegate) {
    application.registerForRemoteNotifications()

    UNUserNotificationCenter.current().delegate = delegate
    UNUserNotificationCenter.current().requestAuthorization(options: [.alert, .badge, .sound]) { success, error in
        guard success else {
            print("Didn't get approval to send push notifications... \(String(describing: error))")
            return
        }
    }
}

extension ApplicationDelegate {
    func application(_ application: Application, didRegisterForRemoteNotificationsWithDeviceToken deviceToken: Data) {
        handleTokenRegistration(deviceToken)
    }

    func application(_ application: Application, didFailToRegisterForRemoteNotificationsWithError error: Error) {
        print("Failed to register for notifications... \(error)")
    }
}

extension ApplicationDelegate: UNUserNotificationCenterDelegate {
    func userNotificationCenter(_ center: UNUserNotificationCenter, willPresent notification: UNNotification, withCompletionHandler completionHandler: @escaping (UNNotificationPresentationOptions) -> Void) {
        completionHandler([[.banner, .sound]])
    }

    func userNotificationCenter(_ center: UNUserNotificationCenter, didReceive response: UNNotificationResponse, withCompletionHandler completionHandler: @escaping () -> Void) {
        completionHandler()
    }
}
