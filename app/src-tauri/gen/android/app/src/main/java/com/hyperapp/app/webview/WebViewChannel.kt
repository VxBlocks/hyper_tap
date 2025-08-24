package com.hyperapp.app.webview

import android.webkit.JavascriptInterface
import android.widget.Toast
import com.hyperapp.app.Logger
import com.hyperapp.app.RustWebView
import kotlinx.coroutines.launch
import kotlinx.serialization.json.Json
import org.json.JSONObject

class WebViewChannel {
  lateinit var webView: RustWebView

  var kotlinToJs: MutableList<WebViewMessage> = mutableListOf()
  var kotlinToJsResults: MutableMap<Int, String> = mutableMapOf()
  var jsToKotlin: MutableList<WebViewMessage> = mutableListOf()
  var jsToKotlinResults: MutableMap<Int, String> = mutableMapOf()

  // 添加一个标志来跟踪事件循环是否已经启动
  private var isEventLoopStarted = false

  // 添加一个锁来确保线程安全
  private val eventLoopLock = Any()

  fun register(webView: RustWebView) {
    this.webView = webView
    webView.addJavascriptInterface(object {
      @JavascriptInterface
      @Suppress("UNUSED")
      fun nextMsg(): String? {
        val msg = kotlinToJs.removeFirstOrNull()
        return if (msg != null) {
          val json = Json.encodeToString(msg)
          Logger.debug("WebViewChannel", "callJs: $json")
          json
        } else {
          null
        }
      }

      @JavascriptInterface
      @Suppress("UNUSED")
      fun addMsg(json: String) {
        val msg = Json.decodeFromString<WebViewMessage>(json)
        Logger.debug("WebViewChannel", "addMsg: $json")
        jsToKotlin.add(msg)
      }

      @JavascriptInterface
      @Suppress("UNUSED")
      fun addResult(id: Int, json: String) {
        Logger.debug("WebViewChannel", "addResult: $id => $json")
        jsToKotlinResults[id] = json
      }
    }, "WebViewChannel")
    startEventLoopIfNeeded()
  }

  suspend fun callJs(method: String, params: JSONObject, callback: ((String) -> Unit)?) {
    // 生成随机ID
    val id = (0..Int.MAX_VALUE).random()

    // 创建消息对象
    val message = WebViewMessage(
      id = id, method = method, params = params.toString()
    )

    // 添加消息到队列
    kotlinToJs.add(message)

    // 通知WebView有新消息
    webView.post {
      webView.evaluateJavascript("window.handleNextMessage()", null)
    }

    if (callback == null) {
      return
    }

    // 等待结果或超时
    val startTime = System.currentTimeMillis()
    val timeout = 5000L // 5秒超时

    while (System.currentTimeMillis() - startTime < timeout) {
      if (kotlinToJsResults.containsKey(id)) {
        val result = kotlinToJsResults.remove(id)
        if (result != null) {
          callback(result)
        }
      }
      // 短暂休眠避免过度占用CPU
      kotlinx.coroutines.delay(10)
    }

    // 超时后报错
    throw RuntimeException("callJs timeout for method: $method with id: $id")
  }

  private fun startEventLoopIfNeeded() {
    synchronized(eventLoopLock) {
      if (!isEventLoopStarted) {
        isEventLoopStarted = true
        kotlinx.coroutines.CoroutineScope(kotlinx.coroutines.Dispatchers.Main).launch {
          startEventLoop()
        }
      }
    }
  }

  suspend fun startEventLoop() {
    while (true) {
      if (jsToKotlin.isNotEmpty()) {
        val message = jsToKotlin.removeFirstOrNull()
        if (message != null) {
          // 处理来自JS的请求
          val result = handleJsRequest(message)

          // 将结果存储到jsToKotlinResults中
          jsToKotlinResults[message.id] = result

          // 通知WebView有结果可用
          webView.post {
            webView.evaluateJavascript("window.handleJsToKotlinResult(${message.id})", null)
          }
        }
      } else {
        // 如果没有消息需要处理，短暂休眠避免过度占用CPU
        kotlinx.coroutines.delay(10)
      }
    }
  }

  private fun handleJsRequest(message: WebViewMessage): String {
    // 这里处理来自JS的请求
    // 根据message.method和message.params执行相应的逻辑
    // 返回处理结果的字符串表示

    // 示例实现 - 你需要根据实际需求修改
    return when (message.method) {
      "toast" -> {
        val jsonObj = JSONObject(message.params)
        val text = jsonObj.getString("text") ?: ""
        webView.post {
          Toast.makeText(webView.context, text, Toast.LENGTH_SHORT).show()
        }
        "success"
      }

      else -> {
        Logger.debug("WebViewChannel", "Unknown method: ${message.method}")
        "unknown method"
      }
    }
  }


  companion object {
    @Volatile
    private var INSTANCE: WebViewChannel? = null

    fun getInstance(): WebViewChannel {
      return INSTANCE ?: synchronized(this) {
        INSTANCE ?: WebViewChannel().also { INSTANCE = it }
      }
    }
  }
}