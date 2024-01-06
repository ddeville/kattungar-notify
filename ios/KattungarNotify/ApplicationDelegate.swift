//
//  ApplicationDelegate.swift
//  Kattungar Notify
//
//  Created by Damien Deville on 1/2/24.
//

#if os(iOS)

import UIKit

class ApplicationDelegate: UIResponder, UIApplicationDelegate, ObservableObject {
    func application(_ application: UIApplication, didFinishLaunchingWithOptions launchOptions: [UIApplication.LaunchOptionsKey: Any]? = nil) -> Bool {
        application.registerForRemoteNotifications()
        requestUserNotificationAuthorization(self)
        return true
    }

    func application(_ application: UIApplication, didRegisterForRemoteNotificationsWithDeviceToken deviceToken: Data) {
        handleTokenRegistration(deviceToken)
    }

    func application(_ application: UIApplication, didFailToRegisterForRemoteNotificationsWithError error: Error) {
        print("Failed to register for notifications... \(error)")
    }
}


#elseif os(macOS)

import Cocoa
import UserNotifications
import SwiftUI

class ApplicationDelegate: NSObject, NSApplicationDelegate, ObservableObject {
    func applicationDidFinishLaunching(_ notification: Notification) {
        let application = notification.object as! NSApplication
        application.registerForRemoteNotifications()
        requestUserNotificationAuthorization(self)
    }

    func application(_ application: NSApplication, didRegisterForRemoteNotificationsWithDeviceToken deviceToken: Data) {
        handleTokenRegistration(deviceToken)
    }
}

#endif

func requestUserNotificationAuthorization(_ delegate: UNUserNotificationCenterDelegate) {
    UNUserNotificationCenter.current().delegate = delegate
    UNUserNotificationCenter.current().requestAuthorization(options: [.alert, .badge, .sound]) { success, error in
        guard success else {
            print("Didn't get approval to send push notifications... \(String(describing: error))")
            return
        }
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
