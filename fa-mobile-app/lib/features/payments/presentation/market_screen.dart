import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'payment_providers.dart';
import '../domain/payment_models.dart';

class MarketScreen extends ConsumerWidget {
  const MarketScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final plansAsync = ref.watch(subscriptionPlansProvider);
    final subscriptionAsync = ref.watch(currentSubscriptionProvider);
    final historyAsync = ref.watch(paymentHistoryProvider);
    final paymentState = ref.watch(paymentNotifierProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Market & Subscriptions',
            style: TextStyle(fontWeight: FontWeight.bold)),
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // 1. Current Subscription Status
            subscriptionAsync.when(
              data: (sub) => _buildStatusBanner(context, sub),
              loading: () => const LinearProgressIndicator(),
              error: (e, s) => const SizedBox.shrink(),
            ),
            const SizedBox(height: 24),

            // 2. Subscription Plans
            const Text('Upgrade Your Plan',
                style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold)),
            const SizedBox(height: 16),
            plansAsync.when(
              data: (plans) => SizedBox(
                height: 420,
                child: ListView.separated(
                  scrollDirection: Axis.horizontal,
                  itemCount: plans.length,
                  separatorBuilder: (_, __) => const SizedBox(width: 16),
                  itemBuilder: (context, index) =>
                      _buildPlanCard(context, ref, plans[index]),
                ),
              ),
              loading: () => const Center(child: CircularProgressIndicator()),
              error: (e, s) => Center(child: Text('Error: $e')),
            ),
            const SizedBox(height: 32),

            // 3. Payment Status Overlay (Only when processing)
            if (paymentState.isLoading ||
                paymentState.status != null ||
                paymentState.error != null)
              _buildPaymentStatus(context, paymentState),

            // 4. Transaction History
            const Text('Payment History',
                style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold)),
            const SizedBox(height: 16),
            historyAsync.when(
              data: (txs) => _buildHistoryList(context, txs),
              loading: () => const Center(child: CircularProgressIndicator()),
              error: (e, s) => const Center(child: Text('No history found')),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildStatusBanner(BuildContext context, UserSubscription sub) {
    final isActive = sub.status == 'active';
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: isActive ? const Color(0xFFECFDF5) : const Color(0xFFFFF7ED),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(
            color: isActive ? Colors.green.shade200 : Colors.orange.shade200),
      ),
      child: Row(
        children: [
          Icon(isActive ? Icons.verified : Icons.info_outline,
              color: isActive ? Colors.green : Colors.orange),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text('Status: ${sub.status.toUpperCase()}',
                    style: TextStyle(
                        fontWeight: FontWeight.bold,
                        color: isActive
                            ? Colors.green.shade900
                            : Colors.orange.shade900)),
                if (sub.expiryDate != null)
                  Text('Expiry: ${sub.expiryDate.toString().split(' ')[0]}',
                      style: const TextStyle(fontSize: 12)),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildPlanCard(
      BuildContext context, WidgetRef ref, SubscriptionPlan plan) {
    final primaryGreen = Theme.of(context).primaryColor;
    return Container(
      width: 280,
      padding: const EdgeInsets.all(24),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(16),
        border: Border.all(color: Colors.grey.shade200, width: 2),
        boxShadow: [
          BoxShadow(
              color: Colors.black.withOpacity(0.05),
              blurRadius: 10,
              offset: const Offset(0, 4))
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(plan.name.toUpperCase(),
              style: TextStyle(
                  letterSpacing: 1.2,
                  color: Colors.grey.shade600,
                  fontWeight: FontWeight.bold,
                  fontSize: 12)),
          const SizedBox(height: 8),
          Text('UGX ${plan.price.toStringAsFixed(0)}',
              style:
                  const TextStyle(fontSize: 28, fontWeight: FontWeight.bold)),
          Text('per ${plan.interval}',
              style: TextStyle(color: Colors.grey.shade500)),
          const SizedBox(height: 24),
          const Divider(),
          const SizedBox(height: 16),
          ...plan.features.map((f) => Padding(
                padding: const EdgeInsets.only(bottom: 12),
                child: Row(
                  children: [
                    const Icon(Icons.check_circle,
                        color: Colors.green, size: 18),
                    const SizedBox(width: 8),
                    Expanded(
                        child: Text(f, style: const TextStyle(fontSize: 14))),
                  ],
                ),
              )),
          const Spacer(),
          SizedBox(
            width: double.infinity,
            child: ElevatedButton(
              style: ElevatedButton.styleFrom(
                  backgroundColor: primaryGreen, foregroundColor: Colors.white),
              onPressed: () => _showPaymentDialog(context, ref, plan),
              child: const Text('Subscribe Now'),
            ),
          )
        ],
      ),
    );
  }

  void _showPaymentDialog(
      BuildContext context, WidgetRef ref, SubscriptionPlan plan) {
    final controller = TextEditingController();
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text('Pay for ${plan.name}'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text('Price: UGX ${plan.price.toStringAsFixed(0)}'),
            const SizedBox(height: 16),
            TextField(
              controller: controller,
              decoration: const InputDecoration(
                  labelText: 'MTN Phone Number', hintText: '25677xxxxxxx'),
              keyboardType: TextInputType.phone,
            ),
          ],
        ),
        actions: [
          TextButton(
              onPressed: () => Navigator.pop(context),
              child: const Text('Cancel')),
          ElevatedButton(
            onPressed: () {
              ref.read(paymentNotifierProvider.notifier).payAndSubscribe(
                  plan.price, controller.text, 'Subscription: ${plan.name}');
              Navigator.pop(context);
            },
            child: const Text('Initiate'),
          ),
        ],
      ),
    );
  }

  Widget _buildPaymentStatus(BuildContext context, PaymentState state) {
    return Container(
      margin: const EdgeInsets.only(bottom: 24),
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
          color: Colors.blue.shade50,
          borderRadius: BorderRadius.circular(12),
          border: Border.all(color: Colors.blue.shade200)),
      child: Row(
        children: [
          if (state.isLoading) const CircularProgressIndicator(strokeWidth: 2),
          if (state.error != null) const Icon(Icons.error, color: Colors.red),
          if (state.status == 'Success!')
            const Icon(Icons.check_circle, color: Colors.green),
          const SizedBox(width: 16),
          Expanded(child: Text(state.error ?? state.status ?? '')),
        ],
      ),
    );
  }

  Widget _buildHistoryList(BuildContext context, List<Transaction> txs) {
    return ListView.builder(
      shrinkWrap: true,
      physics: const NeverScrollableScrollPhysics(),
      itemCount: txs.length,
      itemBuilder: (context, index) {
        final tx = txs[index];
        final isSuccess = tx.status == 'SUCCESSFUL';
        return ListTile(
          contentPadding: EdgeInsets.zero,
          leading: CircleAvatar(
              backgroundColor:
                  isSuccess ? Colors.green.shade50 : Colors.red.shade50,
              child: Icon(isSuccess ? Icons.arrow_downward : Icons.error,
                  color: isSuccess ? Colors.green : Colors.red, size: 18)),
          title: Text(tx.reason,
              style:
                  const TextStyle(fontWeight: FontWeight.bold, fontSize: 14)),
          subtitle: Text(tx.timestamp.toString().split('.')[0]),
          trailing: Text('UGX ${tx.amount.toStringAsFixed(0)}',
              style: const TextStyle(fontWeight: FontWeight.bold)),
        );
      },
    );
  }
}
