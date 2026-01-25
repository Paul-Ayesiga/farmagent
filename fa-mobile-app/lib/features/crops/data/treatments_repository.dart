import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/api/api_client.dart';
import '../domain/treatment.dart';

final treatmentsRepositoryProvider =
    Provider((ref) => TreatmentsRepository(ref.read(apiClientProvider)));

class TreatmentsRepository {
  final ApiClient _client;

  TreatmentsRepository(this._client);

  /// POST /treatments - Create a new treatment
  Future<Treatment?> createTreatment(Map<String, dynamic> treatmentData) async {
    try {
      final response = await _client.post('/treatments', data: treatmentData);
      final data = response.data;
      if (data is Map && data.containsKey('data')) {
        return Treatment.fromJson(data['data']);
      }
      return null;
    } catch (e) {
      return null;
    }
  }

  /// PUT /treatments/:id - Rate treatment effectiveness
  Future<Treatment?> rateTreatment(
      String id, int effectiveness, String? notes) async {
    try {
      final response = await _client.put('/treatments/$id', data: {
        'effectiveness': effectiveness,
        if (notes != null) 'notes': notes,
      });
      final data = response.data;
      if (data is Map && data.containsKey('data')) {
        return Treatment.fromJson(data['data']);
      }
      return null;
    } catch (e) {
      return null;
    }
  }
}
