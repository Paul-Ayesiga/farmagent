import 'package:flutter/foundation.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/api/api_client.dart';
import '../domain/crop.dart';
import '../domain/health_record.dart';
import '../domain/treatment.dart';

final cropsRepositoryProvider =
    Provider((ref) => CropsRepository(ref.read(apiClientProvider)));

class CropsRepository {
  final ApiClient _client;

  CropsRepository(this._client);

  /// GET /crops - List all crops
  Future<List<Crop>> getCrops() async {
    try {
      final response = await _client.get('/crops');
      final data = response.data;
      if (data is Map && data.containsKey('data')) {
        return (data['data'] as List).map((e) => Crop.fromJson(e)).toList();
      }
      return [];
    } catch (e) {
      return [];
    }
  }

  /// POST /crops - Create a new crop
  Future<Crop?> addCrop(Map<String, dynamic> cropData) async {
    try {
      final response = await _client.post('/crops', data: cropData);
      final data = response.data;
      if (data is Map && data.containsKey('data')) {
        return Crop.fromJson(data['data']);
      }
      return null;
    } catch (e) {
      return null;
    }
  }

  /// GET /crops/:id - Get a single crop
  Future<Crop?> getCrop(String id) async {
    try {
      final response = await _client.get('/crops/$id');
      final data = response.data;
      if (data is Map && data.containsKey('data')) {
        return Crop.fromJson(data['data']);
      }
      return null;
    } catch (e) {
      return null;
    }
  }

  /// PUT /crops/:id - Update a crop
  Future<Crop?> updateCrop(String id, Map<String, dynamic> cropData) async {
    try {
      final response = await _client.put('/crops/$id', data: cropData);
      final data = response.data;
      if (data is Map && data.containsKey('data')) {
        return Crop.fromJson(data['data']);
      }
      return null;
    } catch (e) {
      return null;
    }
  }

  /// DELETE /crops/:id - Delete a crop
  Future<bool> deleteCrop(String id) async {
    try {
      await _client.delete('/crops/$id');
      return true;
    } catch (e) {
      return false;
    }
  }

  /// GET /crops/:id/health-records - List health records for a crop
  Future<List<HealthRecord>> getHealthRecords(String cropId) async {
    try {
      final response = await _client.get('/crops/$cropId/health-records');
      final data = response.data;
      debugPrint('Health Records Response: $data');
      if (data is Map && data.containsKey('data')) {
        final list = data['data'] as List;
        debugPrint('Number of health records: ${list.length}');
        return list.map((e) => HealthRecord.fromJson(e)).toList();
      }
      return [];
    } catch (e) {
      debugPrint('Error fetching health records: $e');
      return [];
    }
  }

  /// GET /crops/:id/treatments - List treatments for a crop
  Future<List<Treatment>> getTreatments(String cropId) async {
    try {
      final response = await _client.get('/crops/$cropId/treatments');
      final data = response.data;
      if (data is Map && data.containsKey('data')) {
        return (data['data'] as List)
            .map((e) => Treatment.fromJson(e))
            .toList();
      }
      return [];
    } catch (e) {
      return [];
    }
  }
}
