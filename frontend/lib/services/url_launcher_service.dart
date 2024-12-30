import 'package:flutter/foundation.dart' show kIsWeb;
import 'url_launcher_web.dart';
import 'url_launcher_native.dart';

abstract class UrlLauncherService {
  Future<void> launch(String url);

  // Factory constructor to get the appropriate implementation
  factory UrlLauncherService() {
    if (kIsWeb) {
      return WebUrlLauncher();
    } else {
      return NativeUrlLauncher();
    }
  }
}