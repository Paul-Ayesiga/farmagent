class Crop {
  final String id;
  final String name; // crop_type
  final String variety;
  final String fieldId;
  final String plantingDate;
  final String expectedHarvest;
  final int healthScore;
  final String status;
  final String imageUrl;

  Crop({
    required this.id,
    required this.name,
    required this.variety,
    required this.fieldId,
    required this.plantingDate,
    required this.expectedHarvest,
    required this.healthScore,
    required this.status,
    required this.imageUrl,
  });

  factory Crop.fromJson(Map<String, dynamic> json) {
    final attributes = json['attributes'] ?? json;
    // Relationships handling for JSON:API
    final relationships = json['relationships'] ?? {};
    final fieldData = relationships['field']?['data'] ?? {};
    final fieldId = fieldData['id'] ?? '';

    return Crop(
      id: json['id'] ?? '',
      name: attributes['crop_type'] ?? 'Unknown',
      variety: attributes['variety'] ?? '',
      fieldId: fieldId, // Extracted from relationship
      plantingDate: attributes['planting_date'] ?? '',
      expectedHarvest: attributes['expected_harvest'] ?? '',
      healthScore: attributes['health_score'] ?? 100,
      status: attributes['status'] ?? 'healthy',
      imageUrl: attributes['image_url'] ?? '',
    );
  }

  // To JSON for POST (if needed, though usually POST body is diff structure)
  Map<String, dynamic> toJson() {
    return {
      'crop_type': name,
      'variety': variety,
      'planting_date': plantingDate,
      'expected_harvest': expectedHarvest,
      'field_id': fieldId,
    };
  }

  // Getter for display
  String get fieldLocation =>
      'Field #$fieldId'; // Placeholder until we fetch field name
}
