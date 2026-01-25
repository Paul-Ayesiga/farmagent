import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import '../../../core/api/api_client.dart';
import '../../../core/constants/api_constants.dart';
import '../domain/user.dart';

/// Auth repository provider
final authRepositoryProvider = Provider<AuthRepository>((ref) {
  return AuthRepository(
    apiClient: ref.watch(apiClientProvider),
    storage: ref.watch(secureStorageProvider),
  );
});

/// Repository for authentication operations
class AuthRepository {
  final ApiClient _apiClient;
  final FlutterSecureStorage _storage;

  AuthRepository({
    required ApiClient apiClient,
    required FlutterSecureStorage storage,
  })  : _apiClient = apiClient,
        _storage = storage;

  /// Login with phone/email and password
  Future<LoginResponse> login(String identifier, String password) async {
    final response = await _apiClient.post(
      ApiConstants.login,
      data: {
        'identifier': identifier,
        'password': password,
      },
    );

    final loginResponse = LoginResponse.fromJson(response.data);

    // Store tokens
    await _storage.write(
        key: 'access_token', value: loginResponse.tokens.accessToken);
    await _storage.write(
        key: 'refresh_token', value: loginResponse.tokens.refreshToken);

    return loginResponse;
  }

  /// Register new user
  Future<LoginResponse> register({
    required String firstName,
    required String lastName,
    required String phone,
    required String email,
    required String password,
    String role = 'farmer',
  }) async {
    final response = await _apiClient.post(
      ApiConstants.register,
      data: {
        'first_name': firstName,
        'last_name': lastName,
        'phone': phone,
        'email': email,
        'password': password,
        'role': role,
      },
    );

    final loginResponse = LoginResponse.fromJson(response.data);

    await _storage.write(
        key: 'access_token', value: loginResponse.tokens.accessToken);
    await _storage.write(
        key: 'refresh_token', value: loginResponse.tokens.refreshToken);

    return loginResponse;
  }

  /// Get current user profile
  Future<User> getCurrentUser() async {
    final response = await _apiClient.get(ApiConstants.me);
    return User.fromJson(response.data['user'] ?? response.data);
  }

  /// Logout
  Future<void> logout() async {
    try {
      await _apiClient.post(ApiConstants.logout);
    } finally {
      await _storage.deleteAll();
    }
  }

  /// Forgot Password - Request reset link
  Future<void> forgotPassword(String identifier) async {
    await _apiClient.post(
      '/auth/forgot-password',
      data: {'identifier': identifier},
    );
  }

  /// Reset Password with token
  Future<void> resetPassword(String token, String newPassword) async {
    await _apiClient.post(
      '/auth/reset-password',
      data: {
        'token': token,
        'new_password': newPassword,
      },
    );
  }

  /// Change Password
  Future<void> changePassword(
      String currentPassword, String newPassword) async {
    await _apiClient.post(
      '/auth/change-password',
      data: {
        'current_password': currentPassword,
        'new_password': newPassword,
      },
    );
  }

  /// Verify Email
  Future<void> verifyEmail(String token) async {
    await _apiClient.post(
      '/auth/verify-email',
      data: {'token': token},
    );
  }

  /// Resend Verification Email
  Future<void> resendVerification(String identifier) async {
    await _apiClient.post(
      '/auth/resend-verification',
      data: {'identifier': identifier},
    );
  }

  /// Update Profile
  Future<User> updateProfile(Map<String, dynamic> data) async {
    final response = await _apiClient.put(
      ApiConstants.me,
      data: data,
    );
    return User.fromJson(response.data['user'] ?? response.data);
  }

  /// Check if user is logged in
  Future<bool> isLoggedIn() async {
    final token = await _storage.read(key: 'access_token');
    return token != null && token.isNotEmpty;
  }
}
