import 'package:flutter/material.dart';

class AppTheme {
  static ThemeData get lightTheme {
    // 1. COLOR PALETTE (Strict adherence to guide)
    const primaryGreen = Color(0xFF10B981); // #10b981
    const secondaryBlue = Color(0xFF3B82F6); // #3b82f6
    const accentOrange = Color(0xFFF97316); // #f97316
    const backgroundGray = Color(0xFFF3F4F6); // #f3f4f6
    const textGray = Color(0xFF1F2937); // #1f2937
    const surfaceWhite = Colors.white;
    const errorRed = Color(0xFFEF4444); // #ef4444

    return ThemeData(
      useMaterial3: true,
      scaffoldBackgroundColor: backgroundGray,
      colorScheme: ColorScheme.fromSeed(
        seedColor: primaryGreen,
        primary: primaryGreen,
        secondary: secondaryBlue,
        tertiary: accentOrange,
        error: errorRed,
        surface: surfaceWhite,
        brightness: Brightness.light,
      ),

      // 2. TYPOGRAPHY
      textTheme: const TextTheme(
        headlineLarge: TextStyle(
            color: textGray, fontSize: 32, fontWeight: FontWeight.bold),
        headlineMedium: TextStyle(
            color: textGray, fontSize: 24, fontWeight: FontWeight.bold),
        titleMedium: TextStyle(
            color: textGray, fontSize: 18, fontWeight: FontWeight.w600),
        bodyLarge: TextStyle(color: textGray, fontSize: 16),
        bodyMedium: TextStyle(color: textGray, fontSize: 14),
        labelLarge: TextStyle(
            color: Colors.white, fontSize: 16, fontWeight: FontWeight.w500),
      ),

      appBarTheme: const AppBarTheme(
        backgroundColor: backgroundGray,
        foregroundColor: textGray,
        elevation: 0,
        centerTitle: false,
        surfaceTintColor: Colors.transparent,
      ),

      // 3. COMPONENT STYLES
      // Removed cardTheme to avoid type error; styling handled in widgets or defaults

      elevatedButtonTheme: ElevatedButtonThemeData(
        style: ElevatedButton.styleFrom(
          backgroundColor: primaryGreen,
          foregroundColor: Colors.white,
          elevation: 0,
          minimumSize: const Size(double.infinity, 48),
          padding: const EdgeInsets.symmetric(horizontal: 16),
          shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
          textStyle: const TextStyle(fontSize: 16, fontWeight: FontWeight.w500),
        ),
      ),

      inputDecorationTheme: InputDecorationTheme(
        filled: true,
        fillColor: surfaceWhite,
        contentPadding:
            const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(8),
          borderSide: const BorderSide(color: Color(0xFFD1D5DB)),
        ),
        enabledBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(8),
          borderSide: const BorderSide(color: Color(0xFFD1D5DB)),
        ),
        focusedBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(8),
          borderSide: const BorderSide(color: primaryGreen, width: 2),
        ),
      ),

      // Navigation Bar (Clean Standard)
      navigationBarTheme: NavigationBarThemeData(
        backgroundColor: surfaceWhite,
        height: 64,
        elevation: 8,
        shadowColor: Colors.black.withOpacity(0.1),
        indicatorColor: primaryGreen.withOpacity(0.1),
        labelBehavior: NavigationDestinationLabelBehavior.alwaysShow,
        iconTheme: WidgetStateProperty.resolveWith((states) {
          if (states.contains(WidgetState.selected)) {
            return const IconThemeData(color: primaryGreen, size: 24);
          }
          return const IconThemeData(color: Color(0xFF6B7280), size: 24);
        }),
        labelTextStyle: WidgetStateProperty.resolveWith((states) {
          if (states.contains(WidgetState.selected)) {
            return const TextStyle(
                color: primaryGreen, fontSize: 12, fontWeight: FontWeight.w600);
          }
          return const TextStyle(color: Color(0xFF6B7280), fontSize: 12);
        }),
      ),
    );
  }

  static ThemeData get darkTheme {
    return lightTheme; // Fallback to light until dark mode is spec'd
  }
}
