import 'package:firebase_messaging/firebase_messaging.dart';
import 'notification_service.dart';  // 로컬 알림 서비스

class NotificationHandler {
  // 백그라운드에서 푸시 메시지 수신 처리
  static Future<void> backgroundHandler(RemoteMessage message) async {
    print("Background message: ${message.notification?.title}, ${message.notification?.body}");
    await NotificationService.showNotification(
      message.notification?.title ?? "No Title",
      message.notification?.body ?? "No Body",
    );
  }

  // 포그라운드에서 푸시 메시지 수신 처리
  static Future<void> foregroundHandler(RemoteMessage message) async {
    print("Foreground message: ${message.notification?.title}, ${message.notification?.body}");
    await NotificationService.showNotification(
      message.notification?.title ?? "No Title",
      message.notification?.body ?? "No Body",
    );
  }
}
