import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'screens/login_screen.dart';
import 'screens/home_screen.dart';
import 'services/auth_service.dart';
import 'package:url_strategy/url_strategy.dart';

void main() {
  setPathUrlStrategy();
  runApp(
    ChangeNotifierProvider(
      create: (_) => AuthService(),
      child: const MyApp(),
    ),
  );
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'ISX Portfolio',
      theme: ThemeData(
        primarySwatch: Colors.blue,
        useMaterial3: true,
      ),
      initialRoute: '/',
      onGenerateRoute: (settings) {
        if (settings.name?.startsWith('/auth/callback') == true) {
          final uri = Uri.parse(settings.name!);
          final code = uri.queryParameters['code'];
          if (code != null) {
            context.read<AuthService>().handleAuthCallback(code);
          }
          return MaterialPageRoute(builder: (_) => const HomeScreen());
        }
        
        return MaterialPageRoute(
          builder: (_) => Consumer<AuthService>(
            builder: (context, auth, _) {
              return auth.isAuthenticated ? const HomeScreen() : const LoginScreen();
            },
          ),
        );
      },
    );
  }
}
