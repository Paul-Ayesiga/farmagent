import '../../../core/api/api_client.dart';
import '../domain/payment_models.dart';

class PaymentRepository {
  final ApiClient _client;

  PaymentRepository(this._client);

  /// Fetch available plans
  Future<List<SubscriptionPlan>> getPlans() async {
    final response = await _client.get('/subscriptions/plans');
    final data = response.data['data'] as List;
    return data.map((json) => SubscriptionPlan.fromJson(json)).toList();
  }

  /// Get current subscription status
  Future<UserSubscription> getCurrentSubscription() async {
    final response = await _client.get('/subscriptions');
    return UserSubscription.fromJson(response.data['data']);
  }

  /// Initiate a payment request (MoMo)
  Future<String> initiatePayment(
      double amount, String phone, String reason) async {
    final response = await _client.post('/payments/initiate', data: {
      'amount': amount,
      'phone': phone,
      'reason': reason,
    });
    return response.data['transaction_id'];
  }

  /// Check transaction status
  Future<String> getTransactionStatus(String transactionId) async {
    final response = await _client.get('/payments/$transactionId/status');
    return response.data['status'];
  }

  /// Fetch payment history
  Future<List<Transaction>> getHistory() async {
    final response = await _client.get('/payments/history');
    final data = response.data['data'] as List;
    return data.map((json) => Transaction.fromJson(json)).toList();
  }
}
