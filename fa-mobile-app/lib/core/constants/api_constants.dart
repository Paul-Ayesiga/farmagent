import 'dart:io' show Platform;

/// API Configuration constants
class ApiConstants {
  // Android emulator uses 10.0.2.2 to access host machine's localhost
  // iOS simulator can use localhost directly
  static String get baseUrl {
    const port = '8000';
    if (Platform.isAndroid) {
      return 'http://10.0.2.2:$port/api/v1';
    } else {
      return 'http://localhost:$port/api/v1';
    }
  }

  // Auth endpoints
  static const String login = '/auth/login';
  static const String register = '/auth/register';
  static const String refresh = '/auth/refresh';
  static const String logout = '/auth/logout';
  static const String me = '/auth/me';

  // Crop endpoints
  static const String fields = '/fields';
  static const String crops = '/crops';

  // AI endpoints
  static const String aiChat = '/ai/chat';
  static const String aiAnalyze = '/ai/analyze';
  static const String weather = '/ai/weather';

  // Payment endpoints
  static const String payments = '/payments';
  static const String subscriptions = '/subscriptions';
}
