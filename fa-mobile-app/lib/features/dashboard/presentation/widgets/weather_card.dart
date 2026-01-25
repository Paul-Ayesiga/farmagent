import 'package:flutter/material.dart';
import '../weather_providers.dart';

class WeatherCard extends StatelessWidget {
  final WeatherState weatherState;

  const WeatherCard({super.key, required this.weatherState});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(20),
        border: Border.all(color: Colors.grey.shade200),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.05),
            blurRadius: 20,
            offset: const Offset(0, 10),
          ),
        ],
      ),
      child: switch (weatherState) {
        WeatherLoading() => const SizedBox(
            height: 150,
            child: Center(
              child: CircularProgressIndicator(color: Colors.white),
            ),
          ),
        WeatherError(:final message) => SizedBox(
            height: 150,
            child: Center(
              child: Text(message, style: const TextStyle(color: Colors.white)),
            ),
          ),
        WeatherLoaded(:final weather) => Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        children: [
                          Icon(Icons.location_on,
                              color: Colors.grey.shade600, size: 16),
                          const SizedBox(width: 4),
                          Text(
                            'Kampala, Uganda',
                            style: TextStyle(
                                color: Colors.grey.shade600, fontSize: 14),
                          ),
                        ],
                      ),
                      const SizedBox(height: 8),
                      Text(
                        '${weather.temperature.round()}°C',
                        style: const TextStyle(
                          color: Color(0xFF059669), // Emerald
                          fontSize: 52,
                          fontWeight: FontWeight.bold, // Bolder
                        ),
                      ),
                      Text(
                        weather.description,
                        style: TextStyle(
                            color: Colors.grey.shade800,
                            fontSize: 16,
                            fontWeight: FontWeight.w500),
                      ),
                    ],
                  ),
                  Container(
                    padding: const EdgeInsets.all(16),
                    decoration: BoxDecoration(
                      color: Colors.white.withAlpha(25),
                      borderRadius: BorderRadius.circular(20),
                    ),
                    child: const Icon(Icons.wb_sunny,
                        color: Colors.yellow, size: 48),
                  ),
                ],
              ),
              const SizedBox(height: 20),
              Container(
                padding:
                    const EdgeInsets.symmetric(vertical: 12, horizontal: 16),
                decoration: BoxDecoration(
                  color: Colors.white.withAlpha(25),
                  borderRadius: BorderRadius.circular(12),
                ),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceAround,
                  children: [
                    _WeatherStat(
                      icon: Icons.water_drop,
                      value: '${weather.humidity.round()}%',
                      label: 'Humidity',
                    ),
                    Container(width: 1, height: 30, color: Colors.white24),
                    _WeatherStat(
                      icon: Icons.air,
                      value: '${weather.windSpeed.round()} km/h',
                      label: 'Wind',
                    ),
                  ],
                ),
              ),
            ],
          ),
      },
    );
  }
}

class _WeatherStat extends StatelessWidget {
  final IconData icon;
  final String value;
  final String label;

  const _WeatherStat({
    required this.icon,
    required this.value,
    required this.label,
  });

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        Icon(icon, color: Colors.white70, size: 20),
        const SizedBox(width: 8),
        Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              value,
              style: const TextStyle(
                color: Colors.white,
                fontWeight: FontWeight.bold,
                fontSize: 16,
              ),
            ),
            Text(
              label,
              style: const TextStyle(color: Colors.white70, fontSize: 12),
            ),
          ],
        ),
      ],
    );
  }
}
