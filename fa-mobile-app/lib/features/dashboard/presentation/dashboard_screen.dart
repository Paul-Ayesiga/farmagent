import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:shimmer/shimmer.dart';
import '../../auth/presentation/auth_providers.dart';
import 'weather_providers.dart';
import '../../crops/presentation/crops_providers.dart';

class DashboardScreen extends ConsumerStatefulWidget {
  const DashboardScreen({super.key});

  @override
  ConsumerState<DashboardScreen> createState() => _DashboardScreenState();
}

class _DashboardScreenState extends ConsumerState<DashboardScreen> {
  // No local selectedIndex or bottomNavigationBar

  @override
  Widget build(BuildContext context) {
    final user = ref.watch(currentUserProvider);
    final weatherState = ref.watch(weatherProvider);
    final cropsAsync = ref.watch(cropsProvider);

    final primaryGreen = Theme.of(context).primaryColor;

    return AnnotatedRegion<SystemUiOverlayStyle>(
      value: SystemUiOverlayStyle.dark.copyWith(
        statusBarColor: Colors.transparent,
        statusBarIconBrightness: Brightness.dark,
      ),
      child: Scaffold(
        backgroundColor: const Color(0xFFF3F4F6),
        body: SafeArea(
          child: RefreshIndicator(
            onRefresh: () async {
              ref.invalidate(weatherProvider);
              ref.invalidate(cropsProvider);
            },
            child: SingleChildScrollView(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // 1. Header Section
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            'Welcome back,',
                            style: Theme.of(context).textTheme.bodyMedium,
                          ),
                          Text(
                            user?.displayName ?? 'Farmer',
                            style: Theme.of(context)
                                .textTheme
                                .headlineMedium
                                ?.copyWith(
                                  fontWeight: FontWeight.bold,
                                  fontSize: 24,
                                ),
                          ),
                        ],
                      ),
                      Row(
                        children: [
                          // Weather
                          if (weatherState is WeatherLoaded)
                            Row(
                              children: [
                                const Icon(Icons.wb_sunny_outlined,
                                    color: Colors.orange, size: 24),
                                const SizedBox(width: 4),
                                Text(
                                    '${weatherState.weather.temperature.round()}°',
                                    style: const TextStyle(
                                        fontWeight: FontWeight.bold,
                                        fontSize: 16)),
                              ],
                            ),
                          const SizedBox(width: 16),
                          // Bell
                          Stack(
                            children: [
                              const Icon(Icons.notifications_outlined,
                                  size: 28),
                              Positioned(
                                right: 0,
                                top: 0,
                                child: Container(
                                  width: 10,
                                  height: 10,
                                  decoration: const BoxDecoration(
                                    color: Colors.red,
                                    shape: BoxShape.circle,
                                  ),
                                ),
                              ),
                            ],
                          ),
                          const SizedBox(width: 12),
                          // Language
                          Container(
                            padding: const EdgeInsets.symmetric(
                                horizontal: 8, vertical: 4),
                            decoration: BoxDecoration(
                              border: Border.all(color: Colors.grey.shade300),
                              borderRadius: BorderRadius.circular(16),
                            ),
                            child: const Text('EN',
                                style: TextStyle(
                                    fontWeight: FontWeight.bold, fontSize: 12)),
                          ),
                          const SizedBox(width: 12),
                          // Profile Avatar
                          GestureDetector(
                            onTap: () => context.push('/profile'),
                            child: Container(
                              decoration: BoxDecoration(
                                shape: BoxShape.circle,
                                border: Border.all(
                                  color: primaryGreen,
                                  width: 2,
                                ),
                              ),
                              child: CircleAvatar(
                                radius: 16,
                                backgroundColor: primaryGreen.withAlpha(30),
                                child: Text(
                                  user != null
                                      ? (user.firstName.isNotEmpty
                                          ? user.firstName[0].toUpperCase()
                                          : user.email[0].toUpperCase())
                                      : 'F',
                                  style: TextStyle(
                                    color: primaryGreen,
                                    fontWeight: FontWeight.bold,
                                    fontSize: 14,
                                  ),
                                ),
                              ),
                            ),
                          ),
                        ],
                      ),
                    ],
                  ),
                  const SizedBox(height: 24),

                  // 2. Dynamic Weather & Location Card 🌍
                  if (weatherState is WeatherLoaded)
                    Container(
                      width: double.infinity,
                      padding: const EdgeInsets.all(20),
                      decoration: BoxDecoration(
                        color: Colors.white,
                        borderRadius: BorderRadius.circular(20),
                        boxShadow: [
                          BoxShadow(
                            color: Colors.grey.withOpacity(0.1),
                            blurRadius: 20,
                            offset: const Offset(0, 8),
                          ),
                        ],
                        border: Border.all(color: Colors.grey.shade200),
                      ),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          // Header Row
                          Row(
                            mainAxisAlignment: MainAxisAlignment.spaceBetween,
                            children: [
                              Row(
                                children: [
                                  Icon(Icons.location_on,
                                      color: Colors.grey.shade600, size: 16),
                                  const SizedBox(width: 4),
                                  Text(
                                    weatherState.weather.locationName ??
                                        'Your Location',
                                    style: TextStyle(
                                      color: Colors.grey.shade600,
                                      fontSize: 14,
                                      fontWeight: FontWeight.w500,
                                    ),
                                  ),
                                ],
                              ),
                              IconButton(
                                onPressed: () {
                                  ref.invalidate(weatherProvider);
                                },
                                icon: Icon(Icons.refresh,
                                    color: Colors.grey.shade400, size: 20),
                                tooltip: 'Refresh Weather',
                              ),
                            ],
                          ),
                          const SizedBox(height: 8),

                          // Main Weather Row
                          Row(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              // Temperature
                              Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  Row(
                                    crossAxisAlignment:
                                        CrossAxisAlignment.start,
                                    children: [
                                      Text(
                                        '${weatherState.weather.temperature.round()}',
                                        style: TextStyle(
                                          fontSize: 64,
                                          fontWeight: FontWeight.bold,
                                          color: Colors.grey.shade800,
                                          height: 1,
                                        ),
                                      ),
                                      Padding(
                                        padding: const EdgeInsets.only(top: 8),
                                        child: Text(
                                          '°C',
                                          style: TextStyle(
                                            fontSize: 24,
                                            fontWeight: FontWeight.bold,
                                            color: Colors.grey.shade400,
                                          ),
                                        ),
                                      ),
                                    ],
                                  ),
                                  Text(
                                    weatherState.weather.description.isNotEmpty
                                        ? weatherState.weather.description
                                        : 'Clear skies',
                                    style: TextStyle(
                                      fontSize: 16,
                                      color: Colors.grey.shade600,
                                      fontWeight: FontWeight.w500,
                                    ),
                                  ),
                                ],
                              ),
                              const Spacer(),
                              // Weather Icon
                              Container(
                                padding: const EdgeInsets.all(16),
                                decoration: BoxDecoration(
                                  color: Colors.grey.shade50,
                                  shape: BoxShape.circle,
                                ),
                                child: Icon(
                                  _getWeatherIcon(
                                      weatherState.weather.description),
                                  size: 48,
                                  color:
                                      Colors.orange.shade400, // Make icon pop
                                ),
                              ),
                            ],
                          ),
                          const SizedBox(height: 20),

                          // Weather Details Row
                          Container(
                            padding: const EdgeInsets.all(12),
                            decoration: BoxDecoration(
                              color: Colors.grey.shade50,
                              borderRadius: BorderRadius.circular(12),
                            ),
                            child: Row(
                              mainAxisAlignment: MainAxisAlignment.spaceAround,
                              children: [
                                _buildWeatherDetail(
                                  Icons.water_drop_outlined,
                                  '${weatherState.weather.humidity.round()}%',
                                  'Humidity',
                                  isDark: true,
                                ),
                                Container(
                                  height: 40,
                                  width: 1,
                                  color: Colors.grey.shade200,
                                ),
                                _buildWeatherDetail(
                                  Icons.air,
                                  '${weatherState.weather.windSpeed.round()} km/h',
                                  'Wind',
                                  isDark: true,
                                ),
                                Container(
                                  height: 40,
                                  width: 1,
                                  color: Colors.grey.shade200,
                                ),
                                _buildWeatherDetail(
                                  Icons.thermostat_outlined,
                                  'Feels ${weatherState.weather.temperature.round()}°',
                                  'Real Feel',
                                  isDark: true,
                                ),
                              ],
                            ),
                          ),
                          const SizedBox(height: 16),

                          // Advisories Row
                          Wrap(
                            spacing: 8,
                            runSpacing: 8,
                            children: [
                              // Farming Advisory Badge
                              if (weatherState.weather.farmingAdvisory != null)
                                Container(
                                  padding: const EdgeInsets.symmetric(
                                      horizontal: 12, vertical: 8),
                                  decoration: BoxDecoration(
                                    color: Colors.green.shade400,
                                    borderRadius: BorderRadius.circular(20),
                                  ),
                                  child: Row(
                                    mainAxisSize: MainAxisSize.min,
                                    children: [
                                      const Icon(Icons.agriculture,
                                          color: Colors.white, size: 16),
                                      const SizedBox(width: 6),
                                      Flexible(
                                        child: Text(
                                          weatherState.weather.farmingAdvisory!,
                                          style: const TextStyle(
                                            color: Colors.white,
                                            fontSize: 13,
                                            fontWeight: FontWeight.w500,
                                          ),
                                          maxLines: 1,
                                          overflow: TextOverflow.ellipsis,
                                        ),
                                      ),
                                    ],
                                  ),
                                ),
                              // Spray Advisory Badge
                              if (weatherState.weather.sprayAdvisory != null)
                                Container(
                                  padding: const EdgeInsets.symmetric(
                                      horizontal: 12, vertical: 8),
                                  decoration: BoxDecoration(
                                    color: weatherState.weather.sprayAdvisory!
                                            .toLowerCase()
                                            .contains('good')
                                        ? Colors.teal.shade400
                                        : Colors.orange.shade400,
                                    borderRadius: BorderRadius.circular(20),
                                  ),
                                  child: Row(
                                    mainAxisSize: MainAxisSize.min,
                                    children: [
                                      const Icon(Icons.water,
                                          color: Colors.white, size: 16),
                                      const SizedBox(width: 6),
                                      Flexible(
                                        child: Text(
                                          weatherState.weather.sprayAdvisory!,
                                          style: const TextStyle(
                                            color: Colors.white,
                                            fontSize: 13,
                                            fontWeight: FontWeight.w500,
                                          ),
                                          maxLines: 1,
                                          overflow: TextOverflow.ellipsis,
                                        ),
                                      ),
                                    ],
                                  ),
                                ),
                            ],
                          ),
                        ],
                      ),
                    ),
                  if (weatherState is WeatherError)
                    Container(
                      padding: const EdgeInsets.all(16),
                      decoration: BoxDecoration(
                          color: Colors.red.shade50,
                          borderRadius: BorderRadius.circular(12)),
                      child: Text((weatherState).message,
                          style: const TextStyle(color: Colors.red)),
                    ),

                  const SizedBox(height: 24),

                  // 3. My Crops (Real Data)
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Text('My Crops',
                          style: Theme.of(context)
                              .textTheme
                              .headlineMedium
                              ?.copyWith(fontSize: 20)),
                      TextButton(
                        onPressed: () => context.go('/crops'),
                        child: Text('View All',
                            style: TextStyle(color: primaryGreen)),
                      ),
                    ],
                  ),
                  const SizedBox(height: 12),
                  SizedBox(
                    height: 145,
                    child: cropsAsync.when(
                      data: (crops) {
                        return ListView.separated(
                          scrollDirection: Axis.horizontal,
                          itemCount: crops.length + 1, // +1 for Add button
                          separatorBuilder: (_, __) =>
                              const SizedBox(width: 12),
                          itemBuilder: (context, index) {
                            if (index == crops.length) {
                              return _buildAddCropButton(context);
                            }
                            final crop = crops[index];
                            return _buildCropCard(
                                context,
                                crop.name,
                                crop.fieldLocation,
                                crop.healthScore,
                                crop.status == 'healthy');
                          },
                        );
                      },
                      loading: () => ListView.separated(
                        scrollDirection: Axis.horizontal,
                        itemCount: 3,
                        separatorBuilder: (_, __) => const SizedBox(width: 12),
                        itemBuilder: (_, __) => _buildCropSkeleton(),
                      ),
                      error: (err, stack) => _buildAddCropButton(context),
                    ),
                  ),
                  const SizedBox(height: 24),

                  // 4. Quick Actions
                  Row(
                    children: [
                      Expanded(
                        child: _buildLargeActionButton(
                            context,
                            'Scan Crop',
                            Icons.camera_alt,
                            primaryGreen,
                            () => context.go('/ai/analyze')),
                      ),
                      const SizedBox(width: 16),
                      Expanded(
                        child: _buildLargeActionButton(
                            context,
                            'Market Prices',
                            Icons.show_chart,
                            const Color(0xFFF97316), // Orange 500
                            () => context.go('/payments')),
                      ),
                    ],
                  ),
                  const SizedBox(height: 24),

                  // 5. Recent Activity
                  Text('Recent Activity',
                      style: Theme.of(context)
                          .textTheme
                          .headlineMedium
                          ?.copyWith(fontSize: 20)),
                  const SizedBox(height: 12),
                  _buildActivityItem(
                      'Diagnosis complete',
                      'Maize field - Leaf blight detected',
                      '2h ago',
                      Icons.medical_services_outlined,
                      const Color(0xFFF97316)), // Orange
                  _buildActivityItem('Market Alert', 'Cassava prices up 5%',
                      '5h ago', Icons.trending_up, primaryGreen),
                ],
              ),
            ),
          ),
        ),
        // No Bottom Nav here! Shell handles it.
      ),
    );
  }

  Widget _buildAddCropButton(BuildContext context) {
    return InkWell(
      onTap: () => context.go('/crops'),
      child: Container(
        width: 100,
        decoration: BoxDecoration(
          color: Colors.white,
          borderRadius: BorderRadius.circular(12),
          border:
              Border.all(color: Colors.grey.shade300, style: BorderStyle.solid),
        ),
        child: const Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(Icons.add_circle_outline, color: Colors.grey, size: 32),
            SizedBox(height: 4),
            Text('Add Crop',
                style: TextStyle(color: Colors.grey, fontSize: 12)),
          ],
        ),
      ),
    );
  }

  Widget _buildCropCard(BuildContext context, String name, String field,
      int health, bool healthy) {
    return Container(
      width: 220,
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
        boxShadow: [
          BoxShadow(
              color: Colors.black.withOpacity(0.05),
              blurRadius: 8,
              offset: const Offset(0, 2)),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Row(
            children: [
              Container(
                padding: const EdgeInsets.all(8),
                decoration: BoxDecoration(
                    color: const Color(0xFFECFDF5),
                    borderRadius: BorderRadius.circular(8)), // Green 50
                child: const Icon(Icons.grass,
                    color: Color(0xFF10B981), size: 24), // Green
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(name,
                        style: const TextStyle(
                            fontWeight: FontWeight.bold, fontSize: 16)),
                    Text(field,
                        style:
                            const TextStyle(color: Colors.grey, fontSize: 12)),
                  ],
                ),
              ),
            ],
          ),
          const Spacer(),
          Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  const Text('Health',
                      style: TextStyle(fontSize: 12, color: Colors.grey)),
                  Container(
                    padding:
                        const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                    decoration: BoxDecoration(
                      color: healthy
                          ? const Color(0xFFECFDF5)
                          : const Color(0xFFFFF7ED), // Green 50 vs Orange 50
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Text(healthy ? 'Healthy' : 'Attention',
                        style: TextStyle(
                            fontSize: 10,
                            fontWeight: FontWeight.bold,
                            color: healthy
                                ? const Color(0xFF10B981)
                                : const Color(0xFFF97316))),
                  ),
                ],
              ),
              const SizedBox(height: 8),
              LinearProgressIndicator(
                value: health / 100,
                backgroundColor: Colors.grey.shade200,
                color:
                    healthy ? const Color(0xFF10B981) : const Color(0xFFF97316),
                minHeight: 6,
                borderRadius: BorderRadius.circular(3),
              ),
            ],
          )
        ],
      ),
    );
  }

  Widget _buildLargeActionButton(BuildContext context, String label,
      IconData icon, Color color, VoidCallback onTap) {
    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(12),
      child: Container(
        height: 120,
        decoration: BoxDecoration(
          color: color,
          borderRadius: BorderRadius.circular(12),
          boxShadow: [
            BoxShadow(
                color: color.withOpacity(0.3),
                blurRadius: 8,
                offset: const Offset(0, 4)),
          ],
        ),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(icon, color: Colors.white, size: 40),
            const SizedBox(height: 12),
            Text(label,
                style: const TextStyle(
                    color: Colors.white,
                    fontWeight: FontWeight.bold,
                    fontSize: 16)),
          ],
        ),
      ),
    );
  }

  Widget _buildActivityItem(String title, String subtitle, String time,
      IconData icon, Color iconColor) {
    return Container(
      margin: const EdgeInsets.only(bottom: 12),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
        boxShadow: [
          BoxShadow(
              color: Colors.black.withOpacity(0.05),
              blurRadius: 4,
              offset: const Offset(0, 1))
        ],
      ),
      child: ListTile(
        contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 4),
        leading: Container(
          padding: const EdgeInsets.all(10),
          decoration: BoxDecoration(
            color: iconColor.withOpacity(0.1),
            borderRadius: BorderRadius.circular(8),
          ),
          child: Icon(icon, color: iconColor, size: 24),
        ),
        title: Text(title,
            style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 14)),
        subtitle: Text(subtitle,
            style: const TextStyle(fontSize: 12, color: Colors.grey)),
        trailing: Text(time,
            style: TextStyle(color: Colors.grey.shade400, fontSize: 12)),
      ),
    );
  }

  IconData _getWeatherIcon(String description) {
    final desc = description.toLowerCase();
    if (desc.contains('rain') || desc.contains('drizzle')) {
      return Icons.water_drop;
    } else if (desc.contains('cloud') || desc.contains('overcast')) {
      return Icons.cloud;
    } else if (desc.contains('sun') || desc.contains('clear')) {
      return Icons.wb_sunny;
    } else if (desc.contains('storm') || desc.contains('thunder')) {
      return Icons.thunderstorm;
    } else if (desc.contains('fog') || desc.contains('mist')) {
      return Icons.foggy;
    } else if (desc.contains('wind')) {
      return Icons.air;
    } else if (desc.contains('snow')) {
      return Icons.ac_unit;
    } else {
      return Icons.wb_sunny_outlined;
    }
  }

  Widget _buildWeatherDetail(IconData icon, String value, String label,
      {bool isDark = false}) {
    final textColor = isDark ? Colors.grey.shade800 : Colors.white;
    final subTextColor =
        isDark ? Colors.grey.shade600 : Colors.white.withAlpha(180);
    final iconColor = isDark ? Colors.grey.shade600 : Colors.white70;

    return Column(
      mainAxisSize: MainAxisSize.min,
      children: [
        Icon(icon, color: iconColor, size: 20),
        const SizedBox(height: 4),
        Text(
          value,
          style: TextStyle(
            color: textColor,
            fontSize: 14,
            fontWeight: FontWeight.bold,
          ),
        ),
        Text(
          label,
          style: TextStyle(
            color: subTextColor,
            fontSize: 11,
          ),
        ),
      ],
    );
  }

  Widget _buildCropSkeleton() {
    return Shimmer.fromColors(
      baseColor: Colors.grey.shade300,
      highlightColor: Colors.grey.shade100,
      child: Container(
        width: 220,
        padding: const EdgeInsets.all(12),
        decoration: BoxDecoration(
          color: Colors.white,
          borderRadius: BorderRadius.circular(12),
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            Row(
              children: [
                Container(
                  width: 40,
                  height: 40,
                  decoration: BoxDecoration(
                    color: Colors.white,
                    borderRadius: BorderRadius.circular(8),
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Container(
                        width: 100,
                        height: 14,
                        color: Colors.white,
                      ),
                      const SizedBox(height: 6),
                      Container(
                        width: 60,
                        height: 10,
                        color: Colors.white,
                      ),
                    ],
                  ),
                ),
              ],
            ),
            const Spacer(),
            Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Container(
                      width: 40,
                      height: 10,
                      color: Colors.white,
                    ),
                    Container(
                      width: 50,
                      height: 16,
                      decoration: BoxDecoration(
                        color: Colors.white,
                        borderRadius: BorderRadius.circular(12),
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 8),
                Container(
                  width: double.infinity,
                  height: 6,
                  decoration: BoxDecoration(
                    color: Colors.white,
                    borderRadius: BorderRadius.circular(3),
                  ),
                ),
              ],
            )
          ],
        ),
      ),
    );
  }
}
