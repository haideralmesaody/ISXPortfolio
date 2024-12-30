import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';
import '../config/api_config.dart';
import 'package:url_launcher/url_launcher.dart' as url_launcher;

class AuthService extends ChangeNotifier {
  bool _isAuthenticated = false;
  String? _userEmail;
  String? _userName;

  bool get isAuthenticated => _isAuthenticated;
  String? get userEmail => _userEmail;
  String? get userName => _userName;

  Future<void> initiateGoogleLogin() async {
    try {
      final response = await http.get(Uri.parse('${ApiConfig.baseUrl}/auth/google/login'));
      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        final url = Uri.parse(data['redirect_url']);
        
        if (await url_launcher.canLaunchUrl(url)) {
          await url_launcher.launchUrl(
            url,
            mode: url_launcher.LaunchMode.platformDefault,
          );
        } else {
          print('Could not launch $url');
        }
      }
    } catch (e) {
      print('Error initiating Google login: $e');
    }
  }

  Future<void> handleAuthCallback(String code) async {
    try {
      final response = await http.get(
        Uri.parse('${ApiConfig.baseUrl}/auth/callback?code=$code'),
      );

      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        _userEmail = data['email'];
        _userName = data['name'];
        _isAuthenticated = true;
        notifyListeners();
      } else {
        print('Error in callback: ${response.body}');
      }
    } catch (e) {
      print('Error handling callback: $e');
    }
  }

  void logout() {
    _isAuthenticated = false;
    _userEmail = null;
    _userName = null;
    notifyListeners();
  }
} 