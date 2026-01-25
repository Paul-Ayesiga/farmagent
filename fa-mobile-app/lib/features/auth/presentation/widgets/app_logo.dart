import 'package:flutter/material.dart';

enum AppLogoStyle { vertical, horizontal, iconOnly }

/// Shared app logo widget for consistent branding
class AppLogo extends StatelessWidget {
  final double size;
  final bool showTagline;
  final AppLogoStyle style;

  final Color? color;

  const AppLogo({
    super.key,
    this.size = 80,
    this.showTagline = true,
    this.style = AppLogoStyle.vertical,
    this.color,
  });

  @override
  Widget build(BuildContext context) {
    final iconWidget = Container(
      padding: EdgeInsets.all(size * 0.2),
      decoration: BoxDecoration(
        color: color?.withAlpha(50) ?? Colors.green.shade50,
        shape: BoxShape.circle,
      ),
      child: Icon(
        Icons.agriculture,
        size: size,
        color: color ?? Colors.green.shade700,
      ),
    );

    if (style == AppLogoStyle.iconOnly) {
      return iconWidget;
    }

    if (style == AppLogoStyle.horizontal) {
      return Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          iconWidget,
          const SizedBox(width: 12),
          Text(
            'FarmAgent',
            style: Theme.of(context).textTheme.headlineSmall?.copyWith(
                  color: color ?? Colors.green.shade700,
                  fontWeight: FontWeight.bold,
                ),
          ),
        ],
      );
    }

    // Vertical style (default)
    return Column(
      mainAxisSize: MainAxisSize.min,
      children: [
        iconWidget,
        const SizedBox(height: 16),
        Text(
          'FarmAgent',
          style: Theme.of(context).textTheme.headlineLarge?.copyWith(
                color: Colors.green.shade700,
                fontWeight: FontWeight.bold,
              ),
          textAlign: TextAlign.center,
        ),
        if (showTagline) ...[
          const SizedBox(height: 8),
          Text(
            'Your AI-Powered Farming Assistant',
            style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                  color: Colors.grey.shade600,
                ),
            textAlign: TextAlign.center,
          ),
        ],
      ],
    );
  }
}
