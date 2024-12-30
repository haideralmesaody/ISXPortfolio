import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';  // For ChangeNotifier
import 'package:http/http.dart' as http;
import 'dart:convert';
// ignore: avoid_web_libraries_in_flutter
import 'dart:html' as html;
import '../config/api_config.dart';
import 'logger_service.dart';
import 'url_launcher_service.dart';
import 'navigation_service.dart';

class AuthService extends ChangeNotifier {
  final UrlLauncherService _urlLauncher = UrlLauncherService();
  
  bool _isAuthenticated = false;
  String? _userEmail;
  String? _userName;
  String? _token;

  // Getters
  bool get isAuthenticated => _isAuthenticated;
  String? get userEmail => _userEmail;
  String? get userName => _userName;
  String? get token => _token;

  AuthService() {
    if (kIsWeb) {
      // Listen for messages from popup
      html.window.onMessage.listen((event) {
        if (event.data['type'] == 'oauth-completion') {
          final data = event.data['data'];
          _userEmail = data['email'];
          _userName = data['name'];
          _token = data['token'];
          _isAuthenticated = true;
          notifyListeners();
          NavigationService.navigateTo('/home');
        }
      });
    }
  }

  Future<void> initiateGoogleLogin() async {
    try {
      LoggerService.debug('Starting Google login process');
      final response = await http.get(Uri.parse('${ApiConfig.baseUrl}/auth/google/login'));
      
      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        final url = data['redirect_url'];
        
        LoggerService.debug('Redirecting to Google OAuth URL: $url');
        await _urlLauncher.launch(url);
      } else {
        LoggerService.debug('Failed to get OAuth URL: ${response.statusCode}');
      }
    } catch (e) {
      LoggerService.error('Login error', e);
    }
  }

  Future<bool> handleAuthCallback(String code) async {
    try {
      LoggerService.debug('Handling auth callback with code: ${code.substring(0, 10)}...');
      
      final response = await http.get(
        Uri.parse('${ApiConfig.baseUrl}/auth/callback?code=$code'),
      );

      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        _userEmail = data['email'];
        _userName = data['name'];
        _token = data['token'];
        _isAuthenticated = true;
        
        LoggerService.debug('Auth successful - User: $_userName');
        notifyListeners();
        return true;
      }
      return false;
    } catch (e) {
      LoggerService.error('Auth callback error', e);
      return false;
    }
  }

  void logout() {
    _isAuthenticated = false;
    _userEmail = null;
    _userName = null;
    _token = null;
    notifyListeners();
  }
}