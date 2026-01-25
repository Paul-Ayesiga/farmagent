import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:dropdown_search/dropdown_search.dart';
import '../../../core/utils/app_toast.dart';
import '../domain/field.dart';
import '../data/fields_repository.dart';
import 'crops_providers.dart';

// Common East African crops
const List<String> _commonCrops = [
  'Maize',
  'Beans',
  'Coffee',
  'Bananas',
  'Cassava',
  'Rice',
  'Millet',
  'Sorghum',
  'Groundnuts',
  'Sweet Potatoes',
  'Tomatoes',
  'Cabbage',
  'Onions',
  'Green Peppers',
  'Sugarcane',
  'Tea',
  'Cotton',
  'Sunflower',
  'Soybeans',
  'Irish Potatoes',
];

class AddCropScreen extends ConsumerStatefulWidget {
  const AddCropScreen({super.key});

  @override
  ConsumerState<AddCropScreen> createState() => _AddCropScreenState();
}

class _AddCropScreenState extends ConsumerState<AddCropScreen> {
  final _formKey = GlobalKey<FormState>();
  final _varietyController = TextEditingController();

  String? _selectedCropType;
  String? _selectedFieldId;
  DateTime? _plantingDate;
  DateTime? _harvestDate;
  bool _isLoading = false;

  @override
  Widget build(BuildContext context) {
    final fieldsAsync = ref.watch(fieldsProvider);
    final primaryColor = Theme.of(context).primaryColor;

    return Scaffold(
      appBar: AppBar(title: const Text('Add New Crop')),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16),
        child: Form(
          key: _formKey,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const Text('Crop Details',
                  style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
              const SizedBox(height: 16),

              // 1. Crop Type - Enhanced Dropdown with Search
              DropdownSearch<String>(
                selectedItem: _selectedCropType,
                items: _commonCrops,
                dropdownDecoratorProps: DropDownDecoratorProps(
                  dropdownSearchDecoration: InputDecoration(
                    labelText: 'Crop Type',
                    hintText: 'Select or search for a crop',
                    prefixIcon: Icon(Icons.eco, color: primaryColor),
                    border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(12),
                    ),
                    enabledBorder: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(12),
                      borderSide: BorderSide(color: Colors.grey.shade300),
                    ),
                    focusedBorder: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(12),
                      borderSide: BorderSide(color: primaryColor, width: 2),
                    ),
                    filled: true,
                    fillColor: Colors.grey.shade50,
                  ),
                ),
                popupProps: PopupProps.menu(
                  showSearchBox: true,
                  searchFieldProps: TextFieldProps(
                    decoration: InputDecoration(
                      hintText: 'Search crops...',
                      prefixIcon: const Icon(Icons.search),
                      border: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(12),
                      ),
                    ),
                  ),
                  menuProps: MenuProps(
                    borderRadius: BorderRadius.circular(12),
                    elevation: 8,
                  ),
                  itemBuilder: (context, item, isSelected) =>
                      _buildDropdownItem(
                    item,
                    isSelected,
                    Icons.eco,
                  ),
                ),
                onChanged: (v) => setState(() => _selectedCropType = v),
                validator: (v) =>
                    v == null ? 'Please select a crop type' : null,
              ),
              const SizedBox(height: 16),

              // 2. Variety
              TextFormField(
                controller: _varietyController,
                decoration: InputDecoration(
                  labelText: 'Variety',
                  hintText: 'e.g., Longe 5, NASE 14',
                  prefixIcon: Icon(Icons.category, color: primaryColor),
                  border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(12),
                  ),
                  filled: true,
                  fillColor: Colors.grey.shade50,
                ),
                validator: (v) => v?.isEmpty == true ? 'Required' : null,
              ),
              const SizedBox(height: 16),

              // 3. Field Selection - Enhanced Dropdown
              fieldsAsync.when(
                data: (fields) => DropdownSearch<Field>(
                  selectedItem:
                      fields.where((f) => f.id == _selectedFieldId).firstOrNull,
                  items: fields,
                  itemAsString: (field) =>
                      '${field.name} (${field.sizeAcres} acres)',
                  dropdownDecoratorProps: DropDownDecoratorProps(
                    dropdownSearchDecoration: InputDecoration(
                      labelText: 'Select Field',
                      hintText: 'Choose a field for this crop',
                      prefixIcon: Icon(Icons.map, color: primaryColor),
                      border: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(12),
                      ),
                      enabledBorder: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(12),
                        borderSide: BorderSide(color: Colors.grey.shade300),
                      ),
                      focusedBorder: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(12),
                        borderSide: BorderSide(color: primaryColor, width: 2),
                      ),
                      filled: true,
                      fillColor: Colors.grey.shade50,
                    ),
                  ),
                  popupProps: PopupProps.menu(
                    showSearchBox: fields.length > 5,
                    menuProps: MenuProps(
                      borderRadius: BorderRadius.circular(12),
                      elevation: 8,
                    ),
                    itemBuilder: (context, field, isSelected) => Container(
                      padding: const EdgeInsets.symmetric(
                          horizontal: 16, vertical: 14),
                      decoration: BoxDecoration(
                        color: isSelected
                            ? primaryColor.withAlpha(30)
                            : Colors.transparent,
                        border: Border(
                          bottom: BorderSide(
                              color: Colors.grey.shade200, width: 0.5),
                        ),
                      ),
                      child: Row(
                        children: [
                          Icon(
                            Icons.landscape,
                            color: isSelected
                                ? primaryColor
                                : Colors.grey.shade600,
                            size: 20,
                          ),
                          const SizedBox(width: 12),
                          Expanded(
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Text(
                                  field.name,
                                  style: TextStyle(
                                    fontWeight: isSelected
                                        ? FontWeight.bold
                                        : FontWeight.w500,
                                    color: isSelected
                                        ? primaryColor
                                        : Colors.black87,
                                  ),
                                ),
                                Text(
                                  '${field.sizeAcres} acres • ${field.soilType.toUpperCase()}',
                                  style: TextStyle(
                                    fontSize: 12,
                                    color: Colors.grey.shade600,
                                  ),
                                ),
                              ],
                            ),
                          ),
                          if (isSelected)
                            Icon(Icons.check_circle,
                                color: primaryColor, size: 20),
                        ],
                      ),
                    ),
                  ),
                  onChanged: (field) =>
                      setState(() => _selectedFieldId = field?.id),
                  validator: (v) => v == null ? 'Please select a field' : null,
                ),
                loading: () => const LinearProgressIndicator(),
                error: (e, s) => Text('Error loading fields: $e',
                    style: const TextStyle(color: Colors.red)),
              ),
              const SizedBox(height: 16),

              // 4. Dates
              Row(
                children: [
                  Expanded(
                    child: _buildDatePicker('Planting Date', _plantingDate,
                        (d) => setState(() => _plantingDate = d)),
                  ),
                  const SizedBox(width: 16),
                  Expanded(
                    child: _buildDatePicker('Expected Harvest', _harvestDate,
                        (d) => setState(() => _harvestDate = d)),
                  ),
                ],
              ),
              const SizedBox(height: 32),

              // Submit
              SizedBox(
                width: double.infinity,
                child: ElevatedButton(
                  onPressed: _isLoading ? null : _submit,
                  child: _isLoading
                      ? const SizedBox(
                          height: 20,
                          width: 20,
                          child: CircularProgressIndicator(
                              color: Colors.white, strokeWidth: 2))
                      : const Text('Save Crop'),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildDatePicker(
      String label, DateTime? date, ValueChanged<DateTime> onPicked) {
    return InkWell(
      onTap: () async {
        final picked = await showDatePicker(
          context: context,
          firstDate: DateTime(2020),
          lastDate: DateTime(2030),
          initialDate: date ?? DateTime.now(),
        );
        if (picked != null) onPicked(picked);
      },
      child: InputDecorator(
        decoration: InputDecoration(labelText: label),
        child: Text(
          date != null
              ? "${date.year}-${date.month.toString().padLeft(2, '0')}-${date.day.toString().padLeft(2, '0')}"
              : 'Select Date',
          style: TextStyle(color: date != null ? Colors.black : Colors.grey),
        ),
      ),
    );
  }

  Future<void> _submit() async {
    if (!_formKey.currentState!.validate()) return;
    if (_plantingDate == null || _harvestDate == null) {
      AppToast.warning('Please select both dates');
      return;
    }

    setState(() => _isLoading = true);

    final success = await ref.read(cropsProvider.notifier).addCrop({
      'crop_type': _selectedCropType ?? '',
      'variety': _varietyController.text,
      'field_id': _selectedFieldId,
      'planting_date': _plantingDate!.toIso8601String().split('T')[0],
      'expected_harvest': _harvestDate!.toIso8601String().split('T')[0],
    });

    setState(() => _isLoading = false);

    if (success && mounted) {
      AppToast.success('Crop added successfully!');
      if (context.mounted) {
        context.pop();
      }
    } else if (mounted) {
      AppToast.error('Failed to add crop');
    }
  }

  Widget _buildDropdownItem(String item, bool isSelected, IconData icon) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
      decoration: BoxDecoration(
        color: isSelected
            ? Theme.of(context).primaryColor.withAlpha(30)
            : Colors.transparent,
        border: Border(
          bottom: BorderSide(color: Colors.grey.shade200, width: 0.5),
        ),
      ),
      child: Row(
        children: [
          Icon(
            icon,
            color: isSelected
                ? Theme.of(context).primaryColor
                : Colors.grey.shade600,
            size: 20,
          ),
          const SizedBox(width: 12),
          Text(
            item,
            style: TextStyle(
              fontWeight: isSelected ? FontWeight.bold : FontWeight.normal,
              color:
                  isSelected ? Theme.of(context).primaryColor : Colors.black87,
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
    );
  }
}
