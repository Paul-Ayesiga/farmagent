import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:geolocator/geolocator.dart';
import '../data/weather_repository.dart';
import '../domain/weather.dart';

/// Weather state
sealed class WeatherState {}

class WeatherLoading extends WeatherState {}

class WeatherLoaded extends WeatherState {
  final Weather weather;
  WeatherLoaded(this.weather);
}

class WeatherError extends WeatherState {
  final String message;
  WeatherError(this.message);
}

/// Weather notifier
class WeatherNotifier extends StateNotifier<WeatherState> {
  final WeatherRepository _repository;

  WeatherNotifier(this._repository) : super(WeatherLoading()) {
    fetchWeather(); // Fetch automatically on init
  }

  Future<void> fetchWeather() async {
    state = WeatherLoading();
    try {
      // 1. Check permissions and get location
      final position = await _determinePosition();

      // 2. Fetch using coordinates
      final weather = await _repository.getWeather(
        latitude: position.latitude,
        longitude: position.longitude,
      );
      state = WeatherLoaded(weather);
    } catch (e) {
      if (e is String) {
        state = WeatherError(e);
      } else {
        state = WeatherError(
            "Unable to get weather data. Please check your connection and location permissions.");
      }
    }
  }

  Future<Position> _determinePosition() async {
    bool serviceEnabled;
    LocationPermission permission;

    serviceEnabled = await Geolocator.isLocationServiceEnabled();
    if (!serviceEnabled) {
      return Future.error('Location services are disabled.');
    }

    permission = await Geolocator.checkPermission();
    if (permission == LocationPermission.denied) {
      permission = await Geolocator.requestPermission();
      if (permission == LocationPermission.denied) {
        return Future.error('Location permissions are denied');
      }
    }

    if (permission == LocationPermission.deniedForever) {
      return Future.error(
          'Location permissions are permanently denied, we cannot request permissions.');
    }

    return await Geolocator.getCurrentPosition();
  }
}

/// Weather provider
final weatherProvider =
    StateNotifierProvider<WeatherNotifier, WeatherState>((ref) {
  return WeatherNotifier(ref.watch(weatherRepositoryProvider));
});
