// This is a basic Flutter widget test.
//
// To perform an interaction with a widget in your test, use the WidgetTester
// utility in the flutter_test package. For example, you can send tap and scroll
// gestures. You can also use WidgetTester to find child widgets in the widget
// tree, read text, and verify that the values of widget properties are correct.

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:provider/provider.dart';
import 'package:isxportfolio/main.dart';
import 'package:isxportfolio/services/auth_service.dart';

void main() {
  testWidgets('Basic app test', (WidgetTester tester) async {
    // Build our app and trigger a frame
    await tester.pumpWidget(
      ChangeNotifierProvider(
        create: (_) => AuthService(),
        child: const MyApp(),
      ),
    );

    // Verify that we have a MaterialApp
    expect(find.byType(MaterialApp), findsOneWidget);

    // Verify that we start with the login screen
    expect(find.text('ISX Portfolio'), findsOneWidget);
  });
}
