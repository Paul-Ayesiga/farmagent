import 'dart:io';
import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/api/api_client.dart';

final scanRepositoryProvider =
    Provider((ref) => ScanRepository(ref.read(apiClientProvider)));

class ScanRepository {
  final ApiClient _client;

  ScanRepository(this._client);

  Future<ScanResult> analyzeCrop(File imageFile,
      {String? cropId, String? cropType}) async {
    try {
      final formData = FormData.fromMap({
        'file': await MultipartFile.fromFile(
          imageFile.path,
          filename: imageFile.path.split('/').last,
        ),
        if (cropType != null) 'crop_type': cropType,
        if (cropId != null) 'crop_id': cropId,
      });

      final response = await _client.post('/ai/analyze', data: formData);
      return ScanResult.fromJson(response.data);
    } catch (e) {
      throw Exception('Analysis failed: $e');
    }
  }
}

class TopPrediction {
  final String diseaseClass;
  final String disease;
  final String crop;
  final double confidence;

  TopPrediction({
    required this.diseaseClass,
    required this.disease,
    required this.crop,
    required this.confidence,
  });

  factory TopPrediction.fromJson(Map<String, dynamic> json) {
    return TopPrediction(
      diseaseClass: json['class'] ?? '',
      disease: json['disease'] ?? '',
      crop: json['crop'] ?? 'Unknown',
      confidence: (json['confidence'] ?? 0).toDouble(),
    );
  }
}

class ScanResult {
  final String diseaseClass;
  final String diseaseName;
  final String cropType;
  final double confidence;
  final bool isHealthy;
  final String severity;
  final int healthScore;
  final String message;
  final List<TopPrediction> topPredictions;

  ScanResult({
    required this.diseaseClass,
    required this.diseaseName,
    required this.cropType,
    required this.confidence,
    required this.isHealthy,
    required this.severity,
    required this.healthScore,
    required this.message,
    required this.topPredictions,
  });

  factory ScanResult.fromJson(Map<String, dynamic> json) {
    final predictions = (json['top_predictions'] as List<dynamic>?)
            ?.map((p) => TopPrediction.fromJson(p))
            .toList() ??
        [];

    return ScanResult(
      diseaseClass: json['disease_class'] ?? 'Unknown',
      diseaseName: json['disease_name'] ?? json['disease'] ?? 'Unknown',
      cropType: json['crop_type'] ?? 'Unknown',
      confidence: (json['confidence'] ?? 0).toDouble(),
      isHealthy: json['is_healthy'] ?? false,
      severity: json['severity'] ?? 'unknown',
      healthScore: json['health_score'] ?? 0,
      message: json['message'] ?? '',
      topPredictions: predictions,
    );
  }

  // Legacy getters for compatibility
  String get disease => diseaseName;
  String get treatment =>
      isHealthy ? '' : 'Consult an agronomist for treatment options.';
  String get recommendation => message;
}
