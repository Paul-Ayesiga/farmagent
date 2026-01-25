/// User model
class User {
  final String id;
  final String email;
  final String phone;
  final String firstName;
  final String lastName;
  final String role;
  final bool isVerified;
  final DateTime createdAt;

  User({
    required this.id,
    required this.email,
    required this.phone,
    required this.firstName,
    required this.lastName,
    required this.role,
    required this.isVerified,
    required this.createdAt,
  });

  /// Get display name (first name or phone/email as fallback)
  String get displayName {
    if (firstName.isNotEmpty) return firstName;
    if (lastName.isNotEmpty) return lastName;
    if (phone.isNotEmpty) return phone;
    return email.split('@').first;
  }

  /// Get full name
  String get fullName {
    final parts = [firstName, lastName].where((s) => s.isNotEmpty);
    if (parts.isEmpty) return displayName;
    return parts.join(' ');
  }

  factory User.fromJson(Map<String, dynamic> json) {
    return User(
      id: json['id'] ?? '',
      email: json['email'] ?? '',
      phone: json['phone'] ?? '',
      firstName: json['first_name'] ?? json['firstName'] ?? '',
      lastName: json['last_name'] ?? json['lastName'] ?? '',
      role: json['role'] ?? 'farmer',
      isVerified: json['is_verified'] ?? json['email_verified'] ?? false,
      createdAt: DateTime.tryParse(json['created_at'] ?? '') ?? DateTime.now(),
    );
  }

  Map<String, dynamic> toJson() => {
        'id': id,
        'email': email,
        'phone': phone,
        'first_name': firstName,
        'last_name': lastName,
        'role': role,
        'is_verified': isVerified,
        'created_at': createdAt.toIso8601String(),
      };
}

/// Auth tokens
class AuthTokens {
  final String accessToken;
  final String refreshToken;
  final DateTime expiresAt;

  AuthTokens({
    required this.accessToken,
    required this.refreshToken,
    required this.expiresAt,
  });

  factory AuthTokens.fromJson(Map<String, dynamic> json) {
    // Handle nested tokens structure
    final tokensMap = json['tokens'] ?? json;
    final expiresIn = tokensMap['expires_in'] ?? 3600;
    return AuthTokens(
      accessToken: tokensMap['access_token'] ?? '',
      refreshToken: tokensMap['refresh_token'] ?? '',
      expiresAt: DateTime.now().add(Duration(seconds: expiresIn)),
    );
  }
}

/// Login response
class LoginResponse {
  final User user;
  final AuthTokens tokens;

  LoginResponse({required this.user, required this.tokens});

  factory LoginResponse.fromJson(Map<String, dynamic> json) {
    return LoginResponse(
      user: User.fromJson(json['user'] ?? {}),
      tokens: AuthTokens.fromJson(json),
    );
  }
}
