//
//  Application.swift
//  Kattungar Notify
//
//  Created by Damien Deville on 1/2/24.
//

import SwiftUI

@main
struct KattungarNotifyApp: App {
#if os(iOS)
    @UIApplicationDelegateAdaptor(ApplicationDelegate.self) var delegate: ApplicationDelegate
#elseif os(macOS)
    @NSApplicationDelegateAdaptor(ApplicationDelegate.self) var delegate: ApplicationDelegate
#endif

    var body: some Scene {
        WindowGroup {
            ContentView()
                .environmentObject(delegate)
        }
    }
}
