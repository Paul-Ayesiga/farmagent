import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../data/auth_repository.dart';
import '../domain/user.dart';

/// Auth state
sealed class AuthState {}

class AuthInitial extends AuthState {}

class AuthLoading extends AuthState {}

class AuthAuthenticated extends AuthState {
  final User user;
  AuthAuthenticated(this.user);
}

class AuthUnauthenticated extends AuthState {}

class AuthError extends AuthState {
  final String message;
  AuthError(this.message);
}

/// Auth state notifier
class AuthNotifier extends StateNotifier<AuthState> {
  final AuthRepository _repository;

  AuthNotifier(this._repository) : super(AuthInitial()) {
    _checkAuthStatus();
  }

  Future<void> _checkAuthStatus() async {
    state = AuthLoading();
    try {
      final isLoggedIn = await _repository.isLoggedIn();
      if (isLoggedIn) {
        final user = await _repository.getCurrentUser();
        state = AuthAuthenticated(user);
      } else {
        state = AuthUnauthenticated();
      }
    } catch (e) {
      state = AuthUnauthenticated();
    }
  }

  Future<void> login(String identifier, String password) async {
    state = AuthLoading();
    try {
      final response = await _repository.login(identifier, password);
      state = AuthAuthenticated(response.user);
    } catch (e) {
      state = AuthError(_extractErrorMessage(e));
    }
  }

  Future<void> register({
    required String firstName,
    required String lastName,
    required String phone,
    required String email,
    required String password,
  }) async {
    state = AuthLoading();
    try {
      final response = await _repository.register(
        firstName: firstName,
        lastName: lastName,
        phone: phone,
        email: email,
        password: password,
      );
      state = AuthAuthenticated(response.user);
    } catch (e) {
      state = AuthError(_extractErrorMessage(e));
    }
  }

  /// Refresh the current user data from the API
  Future<void> refreshUser() async {
    try {
      final user = await _repository.getCurrentUser();
      state = AuthAuthenticated(user);
    } catch (e) {
      // Keep current state if refresh fails
    }
  }

  Future<void> logout() async {
    await _repository.logout();
    state = AuthUnauthenticated();
  }

  Future<void> forgotPassword(String identifier) async {
    state = AuthLoading();
    try {
      await _repository.forgotPassword(identifier);
      state =
          AuthUnauthenticated(); // Stay unauthenticated but operation succeeded
    } catch (e) {
      state = AuthError(_extractErrorMessage(e));
    }
  }

  String _extractErrorMessage(dynamic error) {
    final errorStr = error.toString().toLowerCase();
    if (errorStr.contains('invalid credentials')) {
      return 'Invalid phone number or password';
    }
    if (errorStr.contains('already exists')) {
      return 'Account already exists with this phone number';
    }
    if (errorStr.contains('too many requests') || errorStr.contains('429')) {
      return 'Too many attempts. Please wait a moment.';
    }
    return 'Something went wrong. Please try again.';
  }
}

/// Auth notifier provider
final authProvider = StateNotifierProvider<AuthNotifier, AuthState>((ref) {
  return AuthNotifier(ref.watch(authRepositoryProvider));
});

/// Current user provider
final currentUserProvider = Provider<User?>((ref) {
  final authState = ref.watch(authProvider);
  if (authState is AuthAuthenticated) {
    return authState.user;
  }
  return null;
});
