class Treatment {
  final String id;
  final String cropId;
  final String diseaseName;
  final String treatmentName;
  final String treatmentType; // 'organic' or 'chemical'
  final DateTime applicationDate;
  final int? cost;
  final int? effectiveness; // 1-5 rating
  final String? notes;
  final DateTime createdAt;

  Treatment({
    required this.id,
    required this.cropId,
    required this.diseaseName,
    required this.treatmentName,
    required this.treatmentType,
    required this.applicationDate,
    this.cost,
    this.effectiveness,
    this.notes,
    required this.createdAt,
  });

  factory Treatment.fromJson(Map<String, dynamic> json) {
    final attributes = json['attributes'] ?? json;
    return Treatment(
      id: json['id'] ?? '',
      cropId: attributes['crop_id'] ?? '',
      diseaseName: attributes['disease_name'] ?? '',
      treatmentName: attributes['treatment_name'] ?? '',
      treatmentType: attributes['treatment_type'] ?? 'organic',
      applicationDate:
          DateTime.tryParse(attributes['application_date'] ?? '') ??
              DateTime.now(),
      cost: attributes['cost'],
      effectiveness: attributes['effectiveness'],
      notes: attributes['notes'],
      createdAt:
          DateTime.tryParse(attributes['created_at'] ?? '') ?? DateTime.now(),
    );
  }

  bool get isOrganic => treatmentType == 'organic';
}
