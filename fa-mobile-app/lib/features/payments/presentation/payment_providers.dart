import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/api/api_client.dart';
import '../data/payment_repository.dart';
import '../domain/payment_models.dart';
import 'dart:async';

final paymentRepositoryProvider = Provider((ref) {
  return PaymentRepository(ref.read(apiClientProvider));
});

final subscriptionPlansProvider = FutureProvider<List<SubscriptionPlan>>((ref) {
  return ref.read(paymentRepositoryProvider).getPlans();
});

final currentSubscriptionProvider = FutureProvider<UserSubscription>((ref) {
  return ref.read(paymentRepositoryProvider).getCurrentSubscription();
});

final paymentHistoryProvider = FutureProvider<List<Transaction>>((ref) {
  return ref.read(paymentRepositoryProvider).getHistory();
});

// State for active payment initiation
class PaymentState {
  final bool isLoading;
  final String? status;
  final String? error;

  PaymentState({this.isLoading = false, this.status, this.error});
}

class PaymentNotifier extends StateNotifier<PaymentState> {
  final PaymentRepository _repository;
  Timer? _pollingTimer;

  PaymentNotifier(this._repository) : super(PaymentState());

  Future<void> payAndSubscribe(
      double amount, String phone, String reason) async {
    state = PaymentState(isLoading: true, status: 'Initiating...');
    try {
      final txId = await _repository.initiatePayment(amount, phone, reason);
      state = PaymentState(isLoading: true, status: 'Confirm on Phone...');

      // Start polling for status
      _startPolling(txId);
    } catch (e) {
      state = PaymentState(error: e.toString());
    }
  }

  void _startPolling(String txId) {
    _pollingTimer?.cancel();
    _pollingTimer = Timer.periodic(const Duration(seconds: 3), (timer) async {
      try {
        final status = await _repository.getTransactionStatus(txId);
        if (status == 'SUCCESSFUL') {
          timer.cancel();
          state = PaymentState(status: 'Success!');
        } else if (status == 'FAILED') {
          timer.cancel();
          state = PaymentState(error: 'Payment Failed');
        }
      } catch (e) {
        // Continue polling unless it's a critical error
      }
    });
  }

  @override
  void dispose() {
    _pollingTimer?.cancel();
    super.dispose();
  }
}

final paymentNotifierProvider =
    StateNotifierProvider<PaymentNotifier, PaymentState>((ref) {
  return PaymentNotifier(ref.read(paymentRepositoryProvider));
});
