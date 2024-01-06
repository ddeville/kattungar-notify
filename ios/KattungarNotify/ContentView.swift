//
//  ContentView.swift
//  Kattungar Notify
//
//  Created by Damien Deville on 1/2/24.
//

import SwiftUI

struct ContentView: View {
    @EnvironmentObject private var delegate: ApplicationDelegate

    var body: some View {
        if (delegate.hasSetupDeviceKey) {
            VStack {
                Image(systemName: "megaphone.fill")
                    .imageScale(.large)
                    .foregroundStyle(.tint)
                Text("Ready to Receive Notifications!")
            }
            .padding()
        } else {
            InputView()
        }
    }
}

import SwiftUI

struct InputView: View {
    @EnvironmentObject private var delegate: ApplicationDelegate

    @State private var inputText: String = ""
    @State private var submittedText: String = ""

    var body: some View {
        VStack {
            TextField("Enter device key here", text: $inputText)
                .textFieldStyle(RoundedBorderTextFieldStyle())
                .padding()

            Button("Submit") {
                if !inputText.isEmpty {
                    UserDefaults.standard.set(inputText, forKey: DeviceKeyDefaultsKey)
                    delegate.handleApplicationChange()
                }
            }
        }
    }
}

#Preview {
    ContentView()
}
