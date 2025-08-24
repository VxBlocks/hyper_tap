package com.hyperapp.app

import android.Manifest
import android.content.Intent
import android.content.pm.PackageManager
import android.os.Build
import android.os.Bundle
import android.util.Log
import android.view.View
import android.view.ViewGroup
import android.webkit.WebView
import androidx.activity.enableEdgeToEdge
import androidx.activity.result.contract.ActivityResultContracts
import androidx.core.content.ContextCompat
import com.google.android.gms.tasks.OnCompleteListener
import com.google.firebase.messaging.FirebaseMessaging
import com.gyf.immersionbar.ImmersionBar
import com.hyperapp.app.webview.WebViewChannel
import com.hyperapp.app.webview.WebViewUtils
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch
import android.view.WindowManager
import androidx.appcompat.app.AppCompatActivity
import com.hyperapp.app.firebase.MyFirebaseMessagingService


class MainActivity : TauriActivity() {
  override fun onCreate(savedInstanceState: Bundle?) {
    enableEdgeToEdge()
    super.onCreate(savedInstanceState)

    ImmersionBar.with(this)
      .statusBarDarkFont(true, 0.2f)  // 状态栏字体深色
      .navigationBarColor(android.R.color.white)
      .navigationBarDarkIcon(true)
      .fitsSystemWindows(true)        // 适配系统窗口
      .init()
  }

  private lateinit var mWebView: RustWebView

  // private fun setFullscreen() {
  //   WindowCompat.setDecorFitsSystemWindows(window, false)
  //   window.statusBarColor = Color.TRANSPARENT
  //   window.navigationBarColor = Color.TRANSPARENT
  //   WindowCompat.getInsetsController(window, findViewById(android.R.id.content)).let { controller ->
  //     controller.isAppearanceLightStatusBars = false
  //     controller.isAppearanceLightNavigationBars = false
  //   }
  // }

  override fun onWebViewCreate(webView: WebView) {
    super.onWebViewCreate(webView)
    this.mWebView = webView as RustWebView
    WebViewChannel.getInstance().register(webView)

    mWebView.post {
      val layoutParams = mWebView.layoutParams as ViewGroup.MarginLayoutParams
      layoutParams.bottomMargin = getNavigationBarHeight()
      mWebView.layoutParams = layoutParams
    }

    askNotificationPermission()
    FirebaseMessaging.getInstance().token.addOnCompleteListener(OnCompleteListener { task ->
      if (!task.isSuccessful) {
        Log.w("FCM", "Fetching FCM registration token failed", task.exception)
        return@OnCompleteListener
      }

      // Get new FCM registration token
      val token = task.result

      // Log and toast
      Log.d("FCM", token)
      CoroutineScope(Dispatchers.Default).launch {
        WebViewUtils.saveFcmToken(token)
        MyFirebaseMessagingService.subscribeTopic(100)
      }
    })
  }

  override fun onNewIntent(intent: Intent) {
    super.onNewIntent(intent)
    handleNotificationIntent(intent)
  }

  private fun handleNotificationIntent(intent: Intent?) {
    intent?.let {
      val url = it.getBooleanExtra("notification_url", false)
      if (url && ::mWebView.isInitialized) {
        Logger.info("MainActivity", "Handling notification url")
        CoroutineScope(Dispatchers.Default).launch {
          WebViewUtils.gotoNotification()
        }
      }
    }
  }

  // Declare the launcher at the top of your Activity/Fragment:
  private val requestPermissionLauncher = registerForActivityResult(
    ActivityResultContracts.RequestPermission(),
  ) { isGranted: Boolean ->
    if (isGranted) {
      // FCM SDK (and your app) can post notifications.
    } else {
      // TODO: Inform user that that your app will not show notifications.
    }
  }

  private fun askNotificationPermission() {
    // This is only necessary for API level >= 33 (TIRAMISU)
    if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.TIRAMISU) {
      if (ContextCompat.checkSelfPermission(
          this,
          Manifest.permission.POST_NOTIFICATIONS
        ) == PackageManager.PERMISSION_GRANTED
      ) {
        // FCM SDK (and your app) can post notifications.
      } else if (shouldShowRequestPermissionRationale(Manifest.permission.POST_NOTIFICATIONS)) {
        // TODO: display an educational UI explaining to the user the features that will be enabled
        //       by them granting the POST_NOTIFICATION permission. This UI should provide the user
        //       "OK" and "No thanks" buttons. If the user selects "OK," directly request the permission.
        //       If the user selects "No thanks," allow the user to continue without notifications.
      } else {
        // Directly ask for the permission
        requestPermissionLauncher.launch(Manifest.permission.POST_NOTIFICATIONS)
      }
    }
  }

  private fun getNavigationBarHeight(): Int {
    var result = 0
    val resourceId = resources.getIdentifier("navigation_bar_height", "dimen", "android")
    if (resourceId > 0) {
      result = resources.getDimensionPixelSize(resourceId)
    }
    return result
  }

}
