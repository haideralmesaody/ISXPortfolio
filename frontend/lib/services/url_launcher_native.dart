import 'package:flutter/foundation.dart' show kIsWeb;
import 'url_launcher_service.dart';

class NativeUrlLauncher implements UrlLauncherService {
  @override
  Future<void> launch(String url) async {
    if (!kIsWeb) {
      throw UnsupportedError('Native URL launching not implemented');
    }
  }
}