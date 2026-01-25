import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:geocoding/geocoding.dart';
import '../../../core/api/api_client.dart';
import '../../../core/constants/api_constants.dart';
import '../domain/weather.dart';

/// Weather repository provider
final weatherRepositoryProvider = Provider<WeatherRepository>((ref) {
  return WeatherRepository(apiClient: ref.watch(apiClientProvider));
});

class WeatherRepository {
  final ApiClient _apiClient;

  WeatherRepository({required ApiClient apiClient}) : _apiClient = apiClient;

  /// Get weather by coordinates
  Future<Weather> getWeather({double? latitude, double? longitude}) async {
    // 1. Fetch weather from backend
    final response = await _apiClient.get(
      ApiConstants.weather,
      queryParameters: {
        if (latitude != null) 'latitude': latitude,
        if (longitude != null) 'longitude': longitude,
      },
    );

    var weather = Weather.fromJson(response.data);

    // 2. Reverse geocode if coordinates provided to get location name
    if (latitude != null && longitude != null) {
      try {
        final placemarks = await placemarkFromCoordinates(latitude, longitude);
        if (placemarks.isNotEmpty) {
          final place = placemarks.first;
          // visible: Locality (City) -> SubAdmin (District) -> Admin (Region)
          final locationName = place.locality ??
              place.subAdministrativeArea ??
              place.administrativeArea;

          if (locationName != null && locationName.isNotEmpty) {
            weather = weather.copyWith(locationName: locationName);
          }
        }
      } catch (e) {
        // Ignore geocoding errors, fallback to default/null locationName
        print('Reverse geocoding failed: $e');
      }
    }

    return weather;
  }

  /// Get weather by region name
  Future<Weather> getWeatherByRegion(String region) async {
    final response =
        await _apiClient.get('${ApiConstants.weather}/region/$region');
    return Weather.fromJson(response.data);
  }

  /// Check if it's safe to spray
  Future<SpraySafety> checkSpraySafety(
      {double? latitude, double? longitude}) async {
    final response = await _apiClient.get(
      '/ai/spray-check',
      queryParameters: {
        if (latitude != null) 'latitude': latitude,
        if (longitude != null) 'longitude': longitude,
      },
    );
    return SpraySafety.fromJson(response.data);
  }
}
