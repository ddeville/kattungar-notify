//
//  ContentView.swift
//  Kattungar Notify
//
//  Created by Damien Deville on 1/2/24.
//

import SwiftUI

struct ContentView: View {
    @EnvironmentObject private var delegate: ApplicationDelegate
    @Environment(\.scenePhase) private var scenePhase

    var body: some View {
        ScrollView {
            VStack(alignment: .leading, spacing: 16) {
                Image(systemName: "megaphone.fill")
                    .imageScale(.large)
                    .foregroundStyle(.tint)
                Text(delegate.registrationStatus)
                    .textSelection(.enabled)
                Text("Notifications: \(delegate.permissionStatus)")
                    .font(.callout)
                if delegate.hasSetupDeviceKey {
                    Button {
                        delegate.manuallySendTokenToServer()
                    } label: {
                        HStack {
                            if delegate.isSendingToken {
                                ProgressView()
                                    .controlSize(.small)
                            }
                            Text(delegate.isSendingToken ? "Sending Token..." : "Send Token to Server")
                        }
                    }
                    .disabled(delegate.isSendingToken)
                } else {
                    InputView()
                }

                Divider()
                Text("Diagnostics")
                    .font(.headline)
                Text("Includes app and OS versions, server address, and registration results. Your device key and token are omitted.")
                    .font(.caption)
                    .foregroundStyle(.secondary)
                Button("Copy Diagnostics") {
#if os(iOS)
                    UIPasteboard.general.string = delegate.diagnosticReport
#elseif os(macOS)
                    NSPasteboard.general.clearContents()
                    NSPasteboard.general.setString(delegate.diagnosticReport, forType: .string)
#endif
                }
                Text(delegate.diagnosticReport)
                    .font(.caption.monospaced())
                    .textSelection(.enabled)
                    .frame(maxWidth: .infinity, alignment: .leading)
            }
            .padding()
        }
        .onChange(of: scenePhase) { _, phase in
            if phase == .active {
                delegate.refreshNotificationPermission()
            }
        }
    }
}

struct InputView: View {
    @EnvironmentObject private var delegate: ApplicationDelegate

    @State private var inputText: String = ""

    var body: some View {
        VStack {
            TextField("Enter device key here", text: $inputText)
                .textFieldStyle(RoundedBorderTextFieldStyle())
                .padding()

            Button("Submit") {
                let deviceKey = inputText.trimmingCharacters(in: .whitespacesAndNewlines)
                if !deviceKey.isEmpty {
                    UserDefaults.standard.set(deviceKey, forKey: DeviceKeyDefaultsKey)
                    delegate.handleApplicationChange()
                }
            }
        }
    }
}

#Preview {
    ContentView()
}
