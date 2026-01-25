import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/api/api_client.dart';
import '../domain/field.dart';
import '../domain/crop.dart';

final fieldsRepositoryProvider =
    Provider((ref) => FieldsRepository(ref.read(apiClientProvider)));

class FieldsRepository {
  final ApiClient _client;

  FieldsRepository(this._client);

  /// GET /fields - List all fields
  Future<List<Field>> getFields() async {
    try {
      final response = await _client.get('/fields');
      final data = response.data;
      if (data is Map && data.containsKey('data')) {
        return (data['data'] as List).map((e) => Field.fromJson(e)).toList();
      }
      return [];
    } catch (e) {
      return [];
    }
  }

  /// POST /fields - Create a new field
  Future<Field?> createField(Map<String, dynamic> fieldData) async {
    try {
      final response = await _client.post('/fields', data: fieldData);
      final data = response.data;
      if (data is Map && data.containsKey('data')) {
        return Field.fromJson(data['data']);
      }
      return null;
    } catch (e) {
      return null;
    }
  }

  /// GET /fields/:id - Get a single field
  Future<Field?> getField(String id) async {
    try {
      final response = await _client.get('/fields/$id');
      final data = response.data;
      if (data is Map && data.containsKey('data')) {
        return Field.fromJson(data['data']);
      }
      return null;
    } catch (e) {
      return null;
    }
  }

  /// PUT /fields/:id - Update a field
  Future<Field?> updateField(String id, Map<String, dynamic> fieldData) async {
    try {
      final response = await _client.put('/fields/$id', data: fieldData);
      final data = response.data;
      if (data is Map && data.containsKey('data')) {
        return Field.fromJson(data['data']);
      }
      return null;
    } catch (e) {
      return null;
    }
  }

  /// DELETE /fields/:id - Delete a field
  Future<bool> deleteField(String id) async {
    try {
      await _client.delete('/fields/$id');
      return true;
    } catch (e) {
      return false;
    }
  }

  /// GET /fields/:id/crops - List crops belonging to a field
  Future<List<Crop>> getCropsByField(String fieldId) async {
    try {
      final response = await _client.get('/fields/$fieldId/crops');
      final data = response.data;
      if (data is Map && data.containsKey('data')) {
        return (data['data'] as List).map((e) => Crop.fromJson(e)).toList();
      }
      return [];
    } catch (e) {
      return [];
    }
  }
}

final fieldsProvider = FutureProvider<List<Field>>((ref) async {
  final repo = ref.read(fieldsRepositoryProvider);
  return repo.getFields();
});
