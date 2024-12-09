import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:firebase_core/firebase_core.dart';
import 'package:firebase_messaging/firebase_messaging.dart';
import 'notification_service.dart';  // 로컬 알림 서비스 import
import 'notification_handler.dart'; // 메시지 핸들링 서비스 import
import 'firebase_options.dart';  // Firebase 옵션 import
import 'package:http/http.dart' as http;

void main() async {
  WidgetsFlutterBinding.ensureInitialized();

  // Firebase 초기화
  await Firebase.initializeApp(options: DefaultFirebaseOptions.currentPlatform);

  // Firebase Messaging 인스턴스 생성
  FirebaseMessaging messaging = FirebaseMessaging.instance;
  NotificationSettings settings = await messaging.requestPermission(
    alert: true,
    badge: true,
    sound: true,
  );

  // FCM 토큰 확인 및 로그로 출력
  String? token = await messaging.getToken();
  if (token != null) {
    print("FCM 토큰: $token");
    sendTokenToServer(token);
  } else {
    print("FCM 토큰을 가져올 수 없습니다.");
  }

  // Firebase 메시징 초기화
  FirebaseMessaging.onBackgroundMessage(NotificationHandler.backgroundHandler);
  FirebaseMessaging.onMessage.listen(NotificationHandler.foregroundHandler);

  runApp(MyApp());
}

class MyApp extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Flutter Push Notification Example',
      theme: ThemeData(
        primarySwatch: Colors.blue,
      ),
      home: MyHomePage(),
    );
  }
}

class MyHomePage extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text("Flutter Push Notification Example"),
      ),
      body: Center(
        child: Text("Wait for Push Notifications!"),
      ),
    );
  }
}

void sendTokenToServer(String token) async {
  final url = 'http://10.0.2.2:1323/api/v1/SetDevice';  // 서버 API URL

  final response = await http.post(
    Uri.parse(url),
    headers: {"Content-Type": "application/json"},
    body: json.encode({"device_Token": token, "user_id": 1}),
  );

  if (response.statusCode == 200) {
    print("Token sent successfully");
  } else {
    print("Failed to send token");
  }
}