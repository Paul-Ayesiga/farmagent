class HealthRecord {
  final String id;
  final String cropId;
  final int healthScore;
  final String? imageUrl;
  final String? diseaseDetected;
  final double? confidence;
  final String? severity;
  final String? notes;
  final DateTime createdAt;

  HealthRecord({
    required this.id,
    required this.cropId,
    required this.healthScore,
    this.imageUrl,
    this.diseaseDetected,
    this.confidence,
    this.severity,
    this.notes,
    required this.createdAt,
  });

  factory HealthRecord.fromJson(Map<String, dynamic> json) {
    final attributes = json['attributes'] ?? json;
    return HealthRecord(
      id: json['id'] ?? '',
      cropId: attributes['crop_id'] ?? '',
      healthScore: attributes['health_score'] ?? 0,
      imageUrl: attributes['image_url'],
      diseaseDetected: attributes['disease_detected'],
      confidence: (attributes['confidence'] as num?)?.toDouble(),
      severity: attributes['severity'],
      notes: attributes['notes'],
      createdAt:
          DateTime.tryParse(attributes['created_at'] ?? '') ?? DateTime.now(),
    );
  }

  bool get isHealthy => diseaseDetected == null || diseaseDetected!.isEmpty;
}
