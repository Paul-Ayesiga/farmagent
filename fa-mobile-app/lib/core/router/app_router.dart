import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '../../features/auth/presentation/login_screen.dart';
import '../../features/auth/presentation/register_screen.dart';
import '../../features/auth/presentation/forgot_password_screen.dart';
import '../../features/auth/presentation/reset_password_screen.dart';
import '../../features/auth/presentation/verify_email_screen.dart';
import '../../features/dashboard/presentation/dashboard_screen.dart';
import '../../features/home/presentation/scaffold_with_nav_bar.dart';
import '../../features/crops/presentation/crops_list_screen.dart';
import '../../features/crops/presentation/fields_list_screen.dart';
import '../../features/crops/presentation/add_field_screen.dart';
import '../../features/crops/presentation/add_crop_screen.dart';
import '../../features/crops/presentation/crop_detail_screen.dart';
import '../../features/ai/presentation/scan_screen.dart';
import '../../features/ai/presentation/chat_screen.dart';
import '../../features/auth/presentation/profile_screen.dart';
import '../../features/payments/presentation/market_screen.dart';

final rootNavigatorKey = GlobalKey<NavigatorState>();
final _shellNavigatorKey = GlobalKey<NavigatorState>();

final appRouter = GoRouter(
  navigatorKey: rootNavigatorKey,
  initialLocation: '/login',
  observers: [routeObserver], // Add route observer for refresh detection
  routes: [
    // Auth routes (No Shell)
    GoRoute(path: '/login', builder: (context, state) => const LoginScreen()),
    GoRoute(
        path: '/register', builder: (context, state) => const RegisterScreen()),
    GoRoute(
        path: '/forgot-password',
        builder: (context, state) => const ForgotPasswordScreen()),
    GoRoute(
        path: '/reset-password',
        builder: (context, state) => const ResetPasswordScreen()),
    GoRoute(
        path: '/verify-email',
        builder: (context, state) => const VerifyEmailScreen()),

    // Full-screen routes (Outside Shell)
    GoRoute(
      path: '/ai/analyze',
      builder: (context, state) => const ScanScreen(),
    ),
    GoRoute(
      path: '/ai/chat',
      builder: (context, state) => const ChatScreen(),
    ),
    GoRoute(
      path: '/crops/add',
      builder: (context, state) => const AddCropScreen(),
    ),
    GoRoute(
      path: '/fields/add',
      builder: (context, state) => const AddFieldScreen(),
    ),
    GoRoute(
      path: '/crops/:id/scan',
      builder: (context, state) {
        final source = state.uri.queryParameters['source'] ?? 'camera';
        final name = state.uri.queryParameters['name'];
        return ScanScreen(
          cropId: state.pathParameters['id'],
          cropName: name != null ? Uri.decodeComponent(name) : null,
          sourceType: source == 'gallery'
              ? ImageSourceType.gallery
              : ImageSourceType.camera,
        );
      },
    ),
    GoRoute(
      path: '/profile',
      builder: (context, state) => const ProfileScreen(),
    ),

    // Protected Routes with Shell (Bottom Nav)
    ShellRoute(
      navigatorKey: _shellNavigatorKey,
      builder: (context, state, child) {
        return ScaffoldWithNavBar(state: state, child: child);
      },
      routes: [
        GoRoute(
          path: '/dashboard',
          builder: (context, state) => const DashboardScreen(),
        ),
        GoRoute(
          path: '/fields',
          builder: (context, state) => const FieldsListScreen(),
        ),
        GoRoute(
          path: '/crops',
          builder: (context, state) => const CropsListScreen(),
        ),
        GoRoute(
          path: '/crops/:id',
          builder: (context, state) =>
              CropDetailScreen(cropId: state.pathParameters['id']!),
        ),
        GoRoute(
          path: '/payments',
          builder: (context, state) => const MarketScreen(),
        ),
      ],
    ),

    // Root redirect
    GoRoute(
      path: '/',
      redirect: (context, state) => '/login',
    ),
  ],
);
