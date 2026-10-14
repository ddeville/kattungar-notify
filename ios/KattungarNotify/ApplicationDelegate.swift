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
    @Published var registrationStatus = "Enter a device key to register."
    @Published var permissionStatus = "Checking…"
    @Published private(set) var diagnosticEvents: [String] = []

    private var shouldNotifyAfterManualTokenSend = false

    var diagnosticReport: String {
        let version = Bundle.main.object(forInfoDictionaryKey: "CFBundleShortVersionString") as? String ?? "unknown"
        let build = Bundle.main.object(forInfoDictionaryKey: "CFBundleVersion") as? String ?? "unknown"
        return """
        Kattungar Notify \(version) (\(build))
        OS: \(ProcessInfo.processInfo.operatingSystemVersionString)
        Server: \(TokenRegistrationURL.absoluteString)
        Notification permission: \(permissionStatus)
        Status: \(registrationStatus)

        \(diagnosticEvents.joined(separator: "\n"))
        """
    }

    private func recordDiagnostic(_ message: String) {
        diagnosticEvents.append("\(Date().formatted(date: .omitted, time: .standard)): \(message)")
        if diagnosticEvents.count > 50 {
            diagnosticEvents.removeFirst()
        }
    }

    func refreshNotificationPermission() {
        UNUserNotificationCenter.current().getNotificationSettings { settings in
            DispatchQueue.main.async {
                switch settings.authorizationStatus {
                case .notDetermined: self.permissionStatus = "Not requested"
                case .denied: self.permissionStatus = "Denied — enable notifications in Settings"
                case .authorized: self.permissionStatus = "Authorized"
                case .provisional: self.permissionStatus = "Provisional (quiet notifications)"
                case .ephemeral: self.permissionStatus = "Ephemeral"
                @unknown default: self.permissionStatus = "Unknown"
                }
            }
        }
    }

    func handleApplicationChange() {
        UNUserNotificationCenter.current().delegate = self

        // This will switch between the device key setup and regular view
        hasSetupDeviceKey = UserDefaults.standard.string(forKey: DeviceKeyDefaultsKey) != nil
        refreshNotificationPermission()

        if hasSetupDeviceKey {
            requestToken(isManual: false)

            UNUserNotificationCenter.current().requestAuthorization(options: [.alert, .badge, .sound]) { success, error in
                DispatchQueue.main.async {
                    self.recordDiagnostic("Notification permission request: \(success ? "granted" : "not granted")")
                    if let error {
                        self.recordDiagnostic("Permission error: \(error.localizedDescription)")
                    }
                    self.refreshNotificationPermission()
                }
            }
        }
    }

    func manuallySendTokenToServer() {
        refreshNotificationPermission()
        requestToken(isManual: true)
    }

    private func requestToken(isManual: Bool) {
        guard !isSendingToken else {
            return
        }
        guard UserDefaults.standard.string(forKey: DeviceKeyDefaultsKey) != nil else {
            hasSetupDeviceKey = false
            return
        }

        guard let registerForRemoteNotifications else {
            registrationStatus = "Cannot register before the app finishes launching. Try again."
            recordDiagnostic(registrationStatus)
            return
        }

        isSendingToken = true
        shouldNotifyAfterManualTokenSend = isManual
        registrationStatus = "Waiting for a device token from Apple…"
        recordDiagnostic("Requesting a current token from APNs")
        registerForRemoteNotifications()
    }

    @objc func application(_ application: Application, didRegisterForRemoteNotificationsWithDeviceToken deviceToken: Data) {
        let token = deviceToken.map { data in String(format: "%02.2hhx", data) }.joined()
        let deviceKey = UserDefaults.standard.string(forKey: DeviceKeyDefaultsKey)!
        let isManualTokenSend = shouldNotifyAfterManualTokenSend
        recordDiagnostic("Received APNs token (\(deviceToken.count) bytes)")

        sendTokenToServer(deviceKey: deviceKey, token: token, isManualTokenSend: isManualTokenSend)
    }

    private func sendTokenToServer(deviceKey: String, token: String, isManualTokenSend: Bool) {
        registrationStatus = "Uploading token to server…"
        recordDiagnostic("PUT \(TokenRegistrationURL.absoluteString)")
        registerToken(deviceKey: deviceKey, token: token) { result in
            DispatchQueue.main.async {
                self.isSendingToken = false
                self.shouldNotifyAfterManualTokenSend = false

                switch result {
                case .success(let (response, _)):
                    self.registrationStatus = "Token registered with server (HTTP \(response.statusCode))."
                    self.recordDiagnostic(self.registrationStatus)
                    if isManualTokenSend && response.statusCode == 200 {
                        self.showTokenSentNotification()
                    }
                case .failure(let error):
                    self.registrationStatus = "Token upload failed: \(error.localizedDescription)"
                    if case URLSession.HTTPError.serverSideError(let statusCode) = error {
                        self.registrationStatus = "Token upload failed: HTTP \(statusCode)."
                        if statusCode == 401 {
                            self.registrationStatus = "Server rejected the device key (HTTP 401). Check the key and server device record."
                            self.hasSetupDeviceKey = false
                        }
                    } else if case URLSession.HTTPError.transportError(let underlyingError) = error {
                        let networkError = underlyingError as NSError
                        self.registrationStatus = "Token upload failed: \(networkError.localizedDescription) (\(networkError.domain) \(networkError.code))."
                    }
                    self.recordDiagnostic(self.registrationStatus)
                }
            }
        }
    }

    @objc func application(_ application: Application, didFailToRegisterForRemoteNotificationsWithError error: Error) {
        isSendingToken = false
        shouldNotifyAfterManualTokenSend = false
        let registrationError = error as NSError
        registrationStatus = "APNs registration failed: \(registrationError.localizedDescription) (\(registrationError.domain) \(registrationError.code))."
        recordDiagnostic(registrationStatus)
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
