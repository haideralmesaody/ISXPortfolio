import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'screens/login_screen.dart';
import 'screens/home_screen.dart';
import 'services/auth_service.dart';
import 'package:flutter_web_plugins/flutter_web_plugins.dart';
import 'services/navigation_service.dart';
import 'services/logger_service.dart';

void main() {
  setUrlStrategy(PathUrlStrategy());
  LoggerService.init();
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
      navigatorKey: NavigationService.navigatorKey,
      title: 'ISX Portfolio',
      theme: ThemeData(
        primarySwatch: Colors.blue,
        useMaterial3: true,
      ),
      initialRoute: '/',
      builder: (context, child) {
        return Column(
          children: [
            Expanded(
              child: child ?? const SizedBox(),
            ),
          ],
        );
      },
      onUnknownRoute: (settings) {
        LoggerService.debug('Unknown route: ${settings.name}');
        return MaterialPageRoute(
          builder: (_) => const LoginScreen(),
        );
      },
      routes: {
        '/': (context) => const LoginScreen(),
        '/home': (context) => const HomeScreen(),
        '/login': (context) => const LoginScreen(),
      },
      onGenerateRoute: (settings) {
        LoggerService.debug('Generating route for: ${settings.name}');

        if (settings.name?.startsWith('/auth/callback') == true) {
          LoggerService.debug('Processing callback route');
          final uri = Uri.parse(settings.name!);
          final code = uri.queryParameters['code'];
          
          if (code != null) {
            return MaterialPageRoute(
              builder: (_) => Scaffold(
                body: Center(
                  child: FutureBuilder<bool>(
                    future: Provider.of<AuthService>(_, listen: false)
                      .handleAuthCallback(code),
                    builder: (context, snapshot) {
                      if (snapshot.connectionState == ConnectionState.waiting) {
                        return const CircularProgressIndicator();
                      }
                      
                      if (snapshot.hasError) {
                        return Text('Error: ${snapshot.error}');
                      }
                      
                      if (snapshot.data == true) {
                        return const HomeScreen();
                      }
                      
                      return const LoginScreen();
                    },
                  ),
                ),
              ),
            );
          }
        }
        
        LoggerService.debug('Route not matched, using default handling');
        // Define routes
        switch (settings.name) {
          case '/':
            return MaterialPageRoute(
              builder: (_) => Consumer<AuthService>(
                builder: (context, auth, _) {
                  return auth.isAuthenticated ? const HomeScreen() : const LoginScreen();
                },
              ),
            );
          case '/login':
            return MaterialPageRoute(builder: (_) => const LoginScreen());
          case '/home':
            return MaterialPageRoute(builder: (_) => const HomeScreen());
          case '/debug':
            return MaterialPageRoute(
              builder: (_) => Scaffold(
                body: Center(
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Text('Debug Menu'),
                      ElevatedButton(
                        onPressed: () {
                          LoggerService.debug('Testing auth state');
                          final auth = context.read<AuthService>();
                          LoggerService.debug('Is authenticated: ${auth.isAuthenticated}');
                          LoggerService.debug('User name: ${auth.userName}');
                        },
                        child: Text('Check Auth State'),
                      ),
                      ElevatedButton(
                        onPressed: () {
                          LoggerService.debug('Testing navigation');
                          NavigationService.navigateTo('/home');
                        },
                        child: Text('Go to Home'),
                      ),
                    ],
                  ),
                ),
              ),
            );
          default:
            return MaterialPageRoute(builder: (_) => const LoginScreen());
        }
      },
    );
  }
}
