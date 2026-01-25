class Weather {
  final double temperature;
  final double humidity;
  final double windSpeed;
  final String description;
  final String icon;
  final List<String> farmingAdvisories;
  final SpraySafety? spraySafety;
  final String? farmingAdvisory; // Single string from API
  final String? sprayAdvisory; // Single string from API
  final String? locationName; // Added for UI display

  Weather({
    required this.temperature,
    required this.humidity,
    required this.windSpeed,
    required this.description,
    required this.icon,
    this.farmingAdvisories = const [],
    this.spraySafety,
    this.farmingAdvisory,
    this.sprayAdvisory,
    this.locationName,
  });

  factory Weather.fromJson(Map<String, dynamic> json) {
    return Weather(
      temperature: (json['temperature'] as num).toDouble(),
      humidity: (json['humidity'] as num).toDouble(),
      windSpeed: (json['wind_speed'] as num).toDouble(),
      description: json['weather_description'] ?? '',
      icon: '01d', // Placeholder or map weather_code
      farmingAdvisory: json['farming_advisory'],
      sprayAdvisory: json['spray_advisory'],
      locationName: json['location_name'], // Helper if API returns it
    );
  }

  Weather copyWith({
    double? temperature,
    double? humidity,
    double? windSpeed,
    String? description,
    String? icon,
    List<String>? farmingAdvisories,
    SpraySafety? spraySafety,
    String? farmingAdvisory,
    String? sprayAdvisory,
    String? locationName,
  }) {
    return Weather(
      temperature: temperature ?? this.temperature,
      humidity: humidity ?? this.humidity,
      windSpeed: windSpeed ?? this.windSpeed,
      description: description ?? this.description,
      icon: icon ?? this.icon,
      farmingAdvisories: farmingAdvisories ?? this.farmingAdvisories,
      spraySafety: spraySafety ?? this.spraySafety,
      farmingAdvisory: farmingAdvisory ?? this.farmingAdvisory,
      sprayAdvisory: sprayAdvisory ?? this.sprayAdvisory,
      locationName: locationName ?? this.locationName,
    );
  }
}

class SpraySafety {
  final bool isSafe;
  final String reason;
  final String recommendation;

  SpraySafety({
    required this.isSafe,
    required this.reason,
    required this.recommendation,
  });

  factory SpraySafety.fromJson(Map<String, dynamic> json) {
    return SpraySafety(
      isSafe: json['is_safe'] ?? false,
      reason: json['reason'] ?? '',
      recommendation: json['recommendation'] ?? '',
    );
  }
}
