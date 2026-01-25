import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:dropdown_search/dropdown_search.dart';
import 'package:geolocator/geolocator.dart';
import '../../../core/utils/app_toast.dart';
import '../data/fields_repository.dart';

class AddFieldScreen extends ConsumerStatefulWidget {
  const AddFieldScreen({super.key});

  @override
  ConsumerState<AddFieldScreen> createState() => _AddFieldScreenState();
}

class _AddFieldScreenState extends ConsumerState<AddFieldScreen> {
  final _formKey = GlobalKey<FormState>();
  final _nameController = TextEditingController();
  final _sizeController = TextEditingController();
  String _selectedSoilType = 'loam';
  bool _isLoading = false;

  // GPS Location
  double? _latitude;
  double? _longitude;
  bool _isFetchingLocation = false;
  String? _locationError;

  final _soilTypes = ['loam', 'clay', 'sandy', 'silt', 'peat', 'chalk'];

  @override
  void initState() {
    super.initState();
    // Auto-fetch location when screen opens
    _getCurrentLocation();
  }

  @override
  void dispose() {
    _nameController.dispose();
    _sizeController.dispose();
    super.dispose();
  }

  Future<void> _getCurrentLocation() async {
    setState(() {
      _isFetchingLocation = true;
      _locationError = null;
    });

    try {
      // Check if location services are enabled
      bool serviceEnabled = await Geolocator.isLocationServiceEnabled();
      if (!serviceEnabled) {
        setState(() {
          _locationError = 'Location services are disabled. Please enable GPS.';
          _isFetchingLocation = false;
        });
        return;
      }

      // Check permissions
      LocationPermission permission = await Geolocator.checkPermission();
      if (permission == LocationPermission.denied) {
        permission = await Geolocator.requestPermission();
        if (permission == LocationPermission.denied) {
          setState(() {
            _locationError = 'Location permission denied.';
            _isFetchingLocation = false;
          });
          return;
        }
      }

      if (permission == LocationPermission.deniedForever) {
        setState(() {
          _locationError =
              'Location permissions are permanently denied. Please enable in settings.';
          _isFetchingLocation = false;
        });
        return;
      }

      // Get current position
      Position position = await Geolocator.getCurrentPosition(
        desiredAccuracy: LocationAccuracy.high,
      );

      setState(() {
        _latitude = position.latitude;
        _longitude = position.longitude;
        _isFetchingLocation = false;
        _locationError = null;
      });
    } catch (e) {
      setState(() {
        _locationError = 'Failed to get location: $e';
        _isFetchingLocation = false;
      });
    }
  }

  Future<void> _submitField() async {
    if (!_formKey.currentState!.validate()) return;

    if (_latitude == null || _longitude == null) {
      AppToast.warning('Please get your current location first');
      return;
    }

    setState(() => _isLoading = true);

    final fieldData = {
      'name': _nameController.text.trim(),
      'size_acres': double.tryParse(_sizeController.text) ?? 0,
      'latitude': _latitude!,
      'longitude': _longitude!,
      'soil_type': _selectedSoilType,
    };

    final repo = ref.read(fieldsRepositoryProvider);
    final result = await repo.createField(fieldData);

    setState(() => _isLoading = false);

    if (result != null && mounted) {
      ref.invalidate(fieldsProvider);
      AppToast.success('Field "${result.name}" created successfully!');
      context.pop();
    } else if (mounted) {
      AppToast.error('Failed to create field. Please try again.');
    }
  }

  @override
  Widget build(BuildContext context) {
    final primaryGreen = Theme.of(context).primaryColor;

    return Scaffold(
      appBar: AppBar(
        title: const Text('Add New Field'),
        leading: IconButton(
          icon: const Icon(Icons.close),
          onPressed: () => context.pop(),
        ),
      ),
      body: Form(
        key: _formKey,
        child: ListView(
          padding: const EdgeInsets.all(16),
          children: [
            // Field Name
            TextFormField(
              controller: _nameController,
              decoration: const InputDecoration(
                labelText: 'Field Name',
                hintText: 'e.g., Main Farm, North Plot',
                prefixIcon: Icon(Icons.label_outline),
              ),
              validator: (v) => v == null || v.isEmpty ? 'Enter a name' : null,
            ),
            const SizedBox(height: 16),

            // Size
            TextFormField(
              controller: _sizeController,
              keyboardType: TextInputType.number,
              decoration: const InputDecoration(
                labelText: 'Size (Acres)',
                hintText: 'e.g., 2.5',
                prefixIcon: Icon(Icons.straighten),
              ),
              validator: (v) {
                if (v == null || v.isEmpty) return 'Enter size';
                if (double.tryParse(v) == null) return 'Enter a valid number';
                return null;
              },
            ),
            const SizedBox(height: 16),

            // Soil Type - Enhanced Dropdown
            DropdownSearch<String>(
              selectedItem: _selectedSoilType,
              items: _soilTypes,
              dropdownDecoratorProps: DropDownDecoratorProps(
                dropdownSearchDecoration: InputDecoration(
                  labelText: 'Soil Type',
                  prefixIcon:
                      Icon(Icons.grass, color: Theme.of(context).primaryColor),
                  border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(12),
                  ),
                  enabledBorder: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(12),
                    borderSide: BorderSide(color: Colors.grey.shade300),
                  ),
                  focusedBorder: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(12),
                    borderSide: BorderSide(
                        color: Theme.of(context).primaryColor, width: 2),
                  ),
                  filled: true,
                  fillColor: Colors.grey.shade50,
                ),
              ),
              popupProps: PopupProps.menu(
                showSearchBox: false,
                menuProps: MenuProps(
                  borderRadius: BorderRadius.circular(12),
                  elevation: 8,
                ),
                itemBuilder: (context, item, isSelected) => Container(
                  padding:
                      const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
                  decoration: BoxDecoration(
                    color: isSelected
                        ? Theme.of(context).primaryColor.withAlpha(30)
                        : Colors.transparent,
                    border: Border(
                      bottom:
                          BorderSide(color: Colors.grey.shade200, width: 0.5),
                    ),
                  ),
                  child: Row(
                    children: [
                      Icon(
                        _getSoilTypeIcon(item),
                        color: isSelected
                            ? Theme.of(context).primaryColor
                            : Colors.grey.shade600,
                        size: 20,
                      ),
                      const SizedBox(width: 12),
                      Text(
                        item.toUpperCase(),
                        style: TextStyle(
                          fontWeight:
                              isSelected ? FontWeight.bold : FontWeight.normal,
                          color: isSelected
                              ? Theme.of(context).primaryColor
                              : Colors.black87,
                        ),
                      ),
                      const Spacer(),
                      if (isSelected)
                        Icon(
                          Icons.check_circle,
                          color: Theme.of(context).primaryColor,
                          size: 20,
                        ),
                    ],
                  ),
                ),
              ),
              onChanged: (v) => setState(() => _selectedSoilType = v ?? 'loam'),
            ),
            const SizedBox(height: 24),

            // Location Section - GPS Based
            Text('Field Location',
                style: Theme.of(context).textTheme.titleMedium?.copyWith(
                      fontWeight: FontWeight.bold,
                    )),
            const SizedBox(height: 12),

            Container(
              padding: const EdgeInsets.all(16),
              decoration: BoxDecoration(
                color: _locationError != null
                    ? Colors.red.shade50
                    : _latitude != null
                        ? Colors.green.shade50
                        : Colors.grey.shade100,
                borderRadius: BorderRadius.circular(16),
                border: Border.all(
                  color: _locationError != null
                      ? Colors.red.shade200
                      : _latitude != null
                          ? Colors.green.shade200
                          : Colors.grey.shade300,
                ),
              ),
              child: Column(
                children: [
                  Row(
                    children: [
                      Container(
                        padding: const EdgeInsets.all(12),
                        decoration: BoxDecoration(
                          color: _latitude != null
                              ? primaryGreen.withAlpha(30)
                              : Colors.grey.shade200,
                          shape: BoxShape.circle,
                        ),
                        child: _isFetchingLocation
                            ? SizedBox(
                                width: 24,
                                height: 24,
                                child: CircularProgressIndicator(
                                  strokeWidth: 2,
                                  color: primaryGreen,
                                ),
                              )
                            : Icon(
                                _latitude != null
                                    ? Icons.location_on
                                    : Icons.location_off,
                                color: _latitude != null
                                    ? primaryGreen
                                    : Colors.grey.shade600,
                                size: 24,
                              ),
                      ),
                      const SizedBox(width: 16),
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              _isFetchingLocation
                                  ? 'Getting your location...'
                                  : _locationError != null
                                      ? 'Location Error'
                                      : _latitude != null
                                          ? 'Location Captured ✓'
                                          : 'No Location',
                              style: TextStyle(
                                fontWeight: FontWeight.bold,
                                fontSize: 16,
                                color: _locationError != null
                                    ? Colors.red.shade700
                                    : _latitude != null
                                        ? Colors.green.shade700
                                        : Colors.grey.shade700,
                              ),
                            ),
                            const SizedBox(height: 4),
                            if (_locationError != null)
                              Text(
                                _locationError!,
                                style: TextStyle(
                                  fontSize: 13,
                                  color: Colors.red.shade600,
                                ),
                              )
                            else if (_latitude != null && _longitude != null)
                              Text(
                                '${_latitude!.toStringAsFixed(6)}, ${_longitude!.toStringAsFixed(6)}',
                                style: TextStyle(
                                  fontSize: 14,
                                  color: Colors.grey.shade600,
                                  fontFamily: 'monospace',
                                ),
                              )
                            else
                              Text(
                                'Tap refresh to get GPS coordinates',
                                style: TextStyle(
                                  fontSize: 13,
                                  color: Colors.grey.shade500,
                                ),
                              ),
                          ],
                        ),
                      ),
                      if (!_isFetchingLocation)
                        IconButton(
                          onPressed: _getCurrentLocation,
                          icon: Icon(
                            Icons.refresh,
                            color: primaryGreen,
                          ),
                          tooltip: 'Refresh Location',
                        ),
                    ],
                  ),
                  if (_latitude != null && _longitude != null) ...[
                    const SizedBox(height: 12),
                    Container(
                      padding: const EdgeInsets.symmetric(
                          horizontal: 12, vertical: 8),
                      decoration: BoxDecoration(
                        color: Colors.white,
                        borderRadius: BorderRadius.circular(8),
                      ),
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.spaceAround,
                        children: [
                          _buildCoordDisplay('LAT', _latitude!),
                          Container(
                              height: 30,
                              width: 1,
                              color: Colors.grey.shade300),
                          _buildCoordDisplay('LON', _longitude!),
                        ],
                      ),
                    ),
                  ],
                ],
              ),
            ),
            const SizedBox(height: 32),

            // Submit Button
            SizedBox(
              height: 52,
              child: ElevatedButton(
                onPressed: _isLoading ? null : _submitField,
                style: ElevatedButton.styleFrom(
                  backgroundColor: primaryGreen,
                  foregroundColor: Colors.white,
                ),
                child: _isLoading
                    ? const SizedBox(
                        width: 24,
                        height: 24,
                        child: CircularProgressIndicator(
                            color: Colors.white, strokeWidth: 2))
                    : const Text('Create Field',
                        style: TextStyle(
                            fontSize: 16, fontWeight: FontWeight.bold)),
              ),
            ),
          ],
        ),
      ),
    );
  }

  IconData _getSoilTypeIcon(String soilType) {
    switch (soilType.toLowerCase()) {
      case 'loam':
        return Icons.eco;
      case 'clay':
        return Icons.layers;
      case 'sandy':
        return Icons.grain;
      case 'silt':
        return Icons.water_drop;
      case 'peat':
        return Icons.park;
      case 'chalk':
        return Icons.landscape;
      default:
        return Icons.grass;
    }
  }

  Widget _buildCoordDisplay(String label, double value) {
    return Column(
      children: [
        Text(
          label,
          style: TextStyle(
            fontSize: 11,
            color: Colors.grey.shade500,
            fontWeight: FontWeight.bold,
          ),
        ),
        const SizedBox(height: 2),
        Text(
          value.toStringAsFixed(6),
          style: const TextStyle(
            fontSize: 14,
            fontWeight: FontWeight.w500,
            fontFamily: 'monospace',
          ),
        ),
      ],
    );
  }
}
