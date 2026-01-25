class SubscriptionPlan {
  final String id;
  final String name;
  final double price;
  final String interval; // e.g., "monthly", "yearly"
  final List<String> features;

  SubscriptionPlan({
    required this.id,
    required this.name,
    required this.price,
    required this.interval,
    required this.features,
  });

  factory SubscriptionPlan.fromJson(Map<String, dynamic> json) {
    final attributes = json['attributes'] ?? json;
    return SubscriptionPlan(
      id: json['id'] ?? '',
      name: attributes['name'] ?? '',
      price: (attributes['price'] ?? 0).toDouble(),
      interval: attributes['interval'] ?? '',
      features: List<String>.from(attributes['features'] ?? []),
    );
  }
}

class UserSubscription {
  final String planId;
  final String status; // e.g., "active", "expired"
  final DateTime? expiryDate;

  UserSubscription({
    required this.planId,
    required this.status,
    this.expiryDate,
  });

  factory UserSubscription.fromJson(Map<String, dynamic> json) {
    final attributes = json['attributes'] ?? json;
    return UserSubscription(
      planId: attributes['plan_id'] ?? '',
      status: attributes['status'] ?? 'inactive',
      expiryDate: attributes['expiry_date'] != null
          ? DateTime.parse(attributes['expiry_date'])
          : null,
    );
  }
}

class Transaction {
  final String id;
  final double amount;
  final String status;
  final DateTime timestamp;
  final String reason;

  Transaction({
    required this.id,
    required this.amount,
    required this.status,
    required this.timestamp,
    required this.reason,
  });

  factory Transaction.fromJson(Map<String, dynamic> json) {
    final attributes = json['attributes'] ?? json;
    return Transaction(
      id: json['id'] ?? '',
      amount: (attributes['amount'] ?? 0).toDouble(),
      status: attributes['status'] ?? '',
      timestamp: DateTime.parse(
          attributes['created_at'] ?? DateTime.now().toIso8601String()),
      reason: attributes['reason'] ?? '',
    );
  }
}
