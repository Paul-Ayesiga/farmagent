import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/api/api_client.dart';
import '../domain/health_record.dart';

final healthRecordsRepositoryProvider =
    Provider((ref) => HealthRecordsRepository(ref.read(apiClientProvider)));

class HealthRecordsRepository {
  final ApiClient _client;

  HealthRecordsRepository(this._client);

  /// POST /health-records - Create a new health record
  Future<HealthRecord?> createHealthRecord(
      Map<String, dynamic> recordData) async {
    try {
      final response = await _client.post('/health-records', data: recordData);
      final data = response.data;
      if (data is Map && data.containsKey('data')) {
        return HealthRecord.fromJson(data['data']);
      }
      return null;
    } catch (e) {
      return null;
    }
  }
}
