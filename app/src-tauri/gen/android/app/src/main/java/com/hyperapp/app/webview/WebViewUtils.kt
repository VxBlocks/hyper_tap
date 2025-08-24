package com.hyperapp.app.webview

import org.json.JSONObject

object WebViewUtils {
  suspend fun saveFcmToken(token: String) {
    val channel = WebViewChannel.getInstance()
    channel.callJs("saveFcmToken", JSONObject().put("token", token), null)
  }

  suspend fun gotoNotification() {
    val channel = WebViewChannel.getInstance()
    channel.callJs("gotoNotification", JSONObject(), null)
  }
}