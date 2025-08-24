package com.hyperapp.app.webview

import kotlinx.serialization.Serializable
import org.json.JSONObject

@Serializable
data class WebViewMessage(
  val id: Int,
  val method: String,
  val params: String,
)
