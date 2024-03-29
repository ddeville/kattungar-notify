//
//  Service.swift
//  Kattungar Notify
//
//  Created by Damien Deville on 1/5/24.
//

import Foundation

let TokenDefaultsKey = "kattungar-apns-token"
let DeviceKeyDefaultsKey = "kattungar-device-key"

func registerToken(deviceKey: String, token: String, withCompletionHandler completionHandler: @escaping (URLSession.DataTaskResult) -> Void) {
    var request = URLRequest(url: URL(string: "https://notify.home.kattungar.net/device/token")!)
    request.httpMethod = "PUT"
    request.setValue("Bearer \(deviceKey)", forHTTPHeaderField: "Authorization")
    request.setValue("application/json", forHTTPHeaderField: "Content-Type")
    request.httpBody = try! JSONEncoder().encode(["token": token])

    let task = URLSession.shared.dataTask(with: request) { data, response, error in
        if let error = error {
            completionHandler(Result.failure(URLSession.HTTPError.transportError(error)))
            return
        }
        let response = response as! HTTPURLResponse
        guard (200...299).contains(response.statusCode) else {
            completionHandler(Result.failure(URLSession.HTTPError.serverSideError(response.statusCode)))
            return
        }
        completionHandler(Result.success((response, data!)))
    }
    task.resume()
}

extension URLSession {
    enum HTTPError: Error {
        case transportError(Error)
        case serverSideError(Int)
    }

    typealias DataTaskResult = Result<(HTTPURLResponse, Data), Error>

    func dataTask(with request: URLRequest, completionHandler: @escaping (DataTaskResult) -> Void) -> URLSessionDataTask {
        return self.dataTask(with: request) { (data, response, error) in
            if let error = error {
                completionHandler(Result.failure(HTTPError.transportError(error)))
                return
            }
            let response = response as! HTTPURLResponse
            let status = response.statusCode
            guard (200...299).contains(status) else {
                completionHandler(Result.failure(HTTPError.serverSideError(status)))
                return
            }
            completionHandler(Result.success((response, data!)))
        }
    }
}
