import 'package:flutter/material.dart';
import '../../domain/weather.dart';

class SpraySafetyCard extends StatelessWidget {
  final SpraySafety spraySafety;

  const SpraySafetyCard({super.key, required this.spraySafety});

  @override
  Widget build(BuildContext context) {
    final isSafe = spraySafety.isSafe;

    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: isSafe ? Colors.green.shade50 : Colors.red.shade50,
        borderRadius: BorderRadius.circular(16),
        border: Border.all(
          color: isSafe ? Colors.green.shade200 : Colors.red.shade200,
        ),
      ),
      child: Row(
        children: [
          Container(
            padding: const EdgeInsets.all(12),
            decoration: BoxDecoration(
              color: isSafe ? Colors.green.shade100 : Colors.red.shade100,
              borderRadius: BorderRadius.circular(12),
            ),
            child: Icon(
              isSafe ? Icons.check_circle : Icons.warning,
              color: isSafe ? Colors.green.shade700 : Colors.red.shade700,
              size: 28,
            ),
          ),
          const SizedBox(width: 16),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  isSafe ? 'Safe to Spray Today' : 'Not Recommended to Spray',
                  style: Theme.of(context).textTheme.titleSmall?.copyWith(
                        fontWeight: FontWeight.bold,
                        color: isSafe
                            ? Colors.green.shade700
                            : Colors.red.shade700,
                      ),
                ),
                const SizedBox(height: 4),
                Text(
                  spraySafety.recommendation,
                  style: Theme.of(context).textTheme.bodySmall?.copyWith(
                        color: isSafe
                            ? Colors.green.shade600
                            : Colors.red.shade600,
                      ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
