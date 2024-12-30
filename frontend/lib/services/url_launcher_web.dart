// ignore: avoid_web_libraries_in_flutter
import 'dart:html' as html;
import 'package:flutter/foundation.dart' show kIsWeb;
import 'url_launcher_service.dart';
import 'logger_service.dart';

class WebUrlLauncher implements UrlLauncherService {
  @override
  Future<void> launch(String url) async {
    if (kIsWeb) {
      // Open in popup instead of redirecting main window
      final popup = html.window.open(
        url,
        'Google Login',
        'width=500,height=600,menubar=no,toolbar=no,location=no'
      );
      
      if (popup == null) {
        LoggerService.error('Popup was blocked');
        // Fallback to redirect if popup blocked
        html.window.location.assign(url);
      }
    }
  }
}