package com.hyperapp.app.firebase

import android.app.NotificationChannel
import android.app.NotificationManager
import android.app.PendingIntent
import android.content.Intent
import android.os.Build
import androidx.core.app.NotificationCompat
import com.google.firebase.messaging.FirebaseMessaging
import com.google.firebase.messaging.FirebaseMessagingService
import com.google.firebase.messaging.RemoteMessage
import com.hyperapp.app.Logger
import com.hyperapp.app.MainActivity
import com.hyperapp.app.R
import com.hyperapp.app.webview.WebViewUtils
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch

class MyFirebaseMessagingService : FirebaseMessagingService() {
  override fun onMessageReceived(message: RemoteMessage) {
    super.onMessageReceived(message)
    Logger.info("FirebaseMessagingService", message.notification?.body ?: "")
    // 检查消息是否包含通知
    message.notification?.let {
      showNotification(it.title, it.body)
    } ?: run {
      // 处理数据消息，手动创建通知
      val title = message.data["title"]
      val body = message.data["body"]
      showNotification(title, body)
    }
  }

  override fun onNewToken(token: String) {
    super.onNewToken(token)
    CoroutineScope(Dispatchers.Default).launch {
      WebViewUtils.saveFcmToken(token)
    }
    subscribeTopic(100, 0)
  }


  private fun showNotification(title: String?, body: String?) {
    val channelId = "firebase_notifications"
    val notificationId = System.currentTimeMillis().toInt()

    // 创建通知渠道 (Android 8.0+)
    createNotificationChannel(channelId, "Firebase Notifications")

    // 创建启动应用的意图
    val intent = Intent(this, MainActivity::class.java).apply {
      addFlags(Intent.FLAG_ACTIVITY_CLEAR_TOP)
      putExtra("notification_url", true)
    }

    val pendingIntent = PendingIntent.getActivity(
      this,
      0,
      intent,
      PendingIntent.FLAG_IMMUTABLE or PendingIntent.FLAG_UPDATE_CURRENT
    )

    // 构建通知
    val notificationBuilder = NotificationCompat.Builder(this, channelId)
      .setSmallIcon(R.drawable.ic_launcher_foreground) // 需要确保有此图标资源
      .setContentTitle(title ?: "Notification")
      .setContentText(body ?: "")
      .setAutoCancel(true)
      .setContentIntent(pendingIntent)
      .setPriority(NotificationCompat.PRIORITY_HIGH)

    // 显示通知
    val notificationManager = getSystemService(NotificationManager::class.java)
    notificationManager.notify(notificationId, notificationBuilder.build())
  }

  private fun createNotificationChannel(channelId: String, channelName: String) {
    if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
      val channel = NotificationChannel(
        channelId,
        channelName,
        NotificationManager.IMPORTANCE_HIGH
      )
      val notificationManager = getSystemService(NotificationManager::class.java)
      notificationManager.createNotificationChannel(channel)
    }
  }

  companion object {
    fun subscribeTopic(maxRetries: Int = 3, retryCount: Int = 0) {
      FirebaseMessaging.getInstance().subscribeToTopic("watch_address")
        .addOnCompleteListener { task ->
          if (task.isSuccessful) {
            Logger.info(
              "FirebaseMessagingService",
              "Successfully subscribed to topic: watch_address"
            )
          } else {
            if (retryCount < maxRetries) {
              Logger.warn(
                "FirebaseMessagingService",
                "Failed to subscribe to topic: watch_address, retrying... (${retryCount + 1}/$maxRetries) ${task.exception}",
              )
              // 延迟一段时间后重试
              android.os.Handler(android.os.Looper.getMainLooper()).postDelayed({
                subscribeTopic(maxRetries, retryCount + 1)
              }, (1000 * (retryCount + 1)).toLong()) // 逐步增加延迟时间
            } else {
              Logger.error(
                "FirebaseMessagingService",
                "Failed to subscribe to topic: watch_address after $maxRetries attempts",
                task.exception
              )
            }
          }
        }
    }
  }
}