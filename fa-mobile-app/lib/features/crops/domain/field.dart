class Field {
  final String id;
  final String name;
  final double sizeAcres;
  final String soilType;
  final double latitude;
  final double longitude;

  Field({
    required this.id,
    required this.name,
    required this.sizeAcres,
    required this.soilType,
    required this.latitude,
    required this.longitude,
  });

  factory Field.fromJson(Map<String, dynamic> json) {
    final attributes = json['attributes'] ?? json;
    return Field(
      id: json['id'] ?? '',
      name: attributes['name'] ?? 'Unknown Field',
      sizeAcres: (attributes['size_acres'] ?? 0).toDouble(),
      soilType: attributes['soil_type'] ?? '',
      latitude: (attributes['latitude'] ?? 0).toDouble(),
      longitude: (attributes['longitude'] ?? 0).toDouble(),
    );
  }
}
