//
//  Application.swift
//  Kattungar Notify
//
//  Created by Damien Deville on 1/2/24.
//

import SwiftUI

@main
struct KattungarNotifyApp: App {
    @UIApplicationDelegateAdaptor(ApplicationDelegate.self) var delegate: ApplicationDelegate

    var body: some Scene {
        WindowGroup {
            ContentView()
        }
    }
}
