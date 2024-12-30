import 'package:flutter/material.dart';
import 'logger_service.dart';

class NavigationService {
  static final GlobalKey<NavigatorState> navigatorKey = GlobalKey<NavigatorState>();

  static Future<void> navigateTo(String route) async {
    LoggerService.debug('NavigationService - Navigating to: $route');
    try {
      final navigator = navigatorKey.currentState;
      if (navigator != null && navigator.mounted) {
        LoggerService.debug('Navigator is available');
        await navigator.pushReplacementNamed(route);
        LoggerService.debug('Navigation completed');
      } else {
        LoggerService.debug('Navigator not available, scheduling navigation');
        await Future.delayed(Duration(milliseconds: 100));
        if (navigatorKey.currentState?.mounted ?? false) {
          await navigatorKey.currentState?.pushReplacementNamed(route);
        }
      }
    } catch (e) {
      LoggerService.error('Navigation error', e);
    }
  }
} 