import 'dart:io';
import 'package:camera/camera.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:image_picker/image_picker.dart';
import '../../../core/utils/app_toast.dart';
import '../../crops/data/health_records_repository.dart';
import '../../crops/data/treatments_repository.dart';
import '../../crops/presentation/crop_detail_screen.dart';
import '../data/scan_repository.dart';

enum ImageSourceType { camera, gallery }

class ScanScreen extends ConsumerStatefulWidget {
  final String? cropId;
  final String? cropName;
  final ImageSourceType sourceType;

  const ScanScreen({
    super.key,
    this.cropId,
    this.cropName,
    this.sourceType = ImageSourceType.camera,
  });

  @override
  ConsumerState<ScanScreen> createState() => _ScanScreenState();
}

class _ScanScreenState extends ConsumerState<ScanScreen>
    with WidgetsBindingObserver {
  CameraController? _controller;
  Future<void>? _initializeControllerFuture;
  bool _isProcessing = false;
  File? _capturedImage;
  String? _errorMessage;
  final ImagePicker _picker = ImagePicker();

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addObserver(this);

    if (widget.sourceType == ImageSourceType.gallery) {
      _pickFromGallery();
    } else {
      _initCamera();
    }
  }

  Future<void> _initCamera() async {
    try {
      final cameras = await availableCameras();
      if (cameras.isEmpty) {
        setState(() => _errorMessage = 'No cameras found (Simulator?)');
        return;
      }

      final firstCamera = cameras.first;

      _controller = CameraController(
        firstCamera,
        ResolutionPreset.medium,
        enableAudio: false,
      );

      _initializeControllerFuture = _controller!.initialize();
      await _initializeControllerFuture;
      if (mounted) setState(() {});
    } catch (e) {
      if (mounted) {
        setState(() => _errorMessage = 'Camera Error: $e');
      }
    }
  }

  Future<void> _pickFromGallery() async {
    try {
      final XFile? image = await _picker.pickImage(
        source: ImageSource.gallery,
        maxWidth: 1200,
        maxHeight: 1200,
        imageQuality: 85,
      );

      if (image != null) {
        setState(() {
          _capturedImage = File(image.path);
          _isProcessing = true;
        });

        await _analyzeImage();
      } else {
        if (mounted) {
          if (context.canPop()) {
            context.pop();
          } else {
            context.go('/crops');
          }
        }
      }
    } catch (e) {
      if (mounted) {
        AppToast.error('Failed to pick image: $e');
        context.pop();
      }
    }
  }

  @override
  void dispose() {
    WidgetsBinding.instance.removeObserver(this);
    _controller?.dispose();
    super.dispose();
  }

  @override
  void didChangeAppLifecycleState(AppLifecycleState state) {
    final CameraController? cameraController = _controller;

    if (cameraController == null || !cameraController.value.isInitialized) {
      return;
    }

    if (state == AppLifecycleState.inactive) {
      cameraController.dispose();
    } else if (state == AppLifecycleState.resumed) {
      if (widget.sourceType == ImageSourceType.camera) {
        _initCamera();
      }
    }
  }

  Future<void> _takePicture() async {
    if (_controller == null || !_controller!.value.isInitialized) return;

    try {
      final image = await _controller!.takePicture();
      setState(() {
        _capturedImage = File(image.path);
        _isProcessing = true;
      });

      await _analyzeImage();
    } catch (e) {
      if (mounted) {
        setState(() => _isProcessing = false);
        AppToast.error('Error capturing image: $e');
      }
    }
  }

  Future<void> _analyzeImage() async {
    if (_capturedImage == null) return;

    try {
      final result = await ref.read(scanRepositoryProvider).analyzeCrop(
            _capturedImage!,
            cropId: widget.cropId,
            cropType: widget.cropName,
          );

      setState(() {
        _isProcessing = false;
      });

      if (mounted) _showResultDialog(result);
    } catch (e) {
      if (mounted) {
        setState(() => _isProcessing = false);
        AppToast.error('Analysis failed: $e');
      }
    }
  }

  void _showResultDialog(ScanResult result) {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.white,
      shape: const RoundedRectangleBorder(
          borderRadius: BorderRadius.vertical(top: Radius.circular(24))),
      builder: (ctx) => DraggableScrollableSheet(
        initialChildSize: 0.85,
        minChildSize: 0.5,
        maxChildSize: 0.95,
        expand: false,
        builder: (_, scrollController) => SingleChildScrollView(
          controller: scrollController,
          padding: const EdgeInsets.all(24),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Handle
              Center(
                  child: Container(
                      width: 40,
                      height: 4,
                      decoration: BoxDecoration(
                          color: Colors.grey.shade300,
                          borderRadius: BorderRadius.circular(2)))),
              const SizedBox(height: 20),

              // Image Preview and Crop Info
              if (_capturedImage != null)
                Row(
                  children: [
                    ClipRRect(
                      borderRadius: BorderRadius.circular(12),
                      child: Image.file(_capturedImage!,
                          width: 80, height: 80, fit: BoxFit.cover),
                    ),
                    const SizedBox(width: 16),
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            widget.cropName ?? 'Crop Analysis',
                            style: const TextStyle(
                                fontSize: 18, fontWeight: FontWeight.bold),
                          ),
                          const SizedBox(height: 4),
                          Text(
                            'Detected: ${result.cropType}',
                            style: TextStyle(color: Colors.grey.shade600),
                          ),
                        ],
                      ),
                    ),
                  ],
                ),
              const SizedBox(height: 20),

              // Health Score Card
              Container(
                width: double.infinity,
                padding: const EdgeInsets.all(20),
                decoration: BoxDecoration(
                  gradient: LinearGradient(
                    colors: result.isHealthy
                        ? [Colors.green.shade400, Colors.green.shade600]
                        : [Colors.orange.shade400, Colors.red.shade500],
                  ),
                  borderRadius: BorderRadius.circular(16),
                ),
                child: Column(
                  children: [
                    Icon(
                      result.isHealthy ? Icons.check_circle : Icons.warning,
                      color: Colors.white,
                      size: 48,
                    ),
                    const SizedBox(height: 12),
                    Text(
                      result.isHealthy ? 'Healthy' : result.diseaseName,
                      style: const TextStyle(
                          color: Colors.white,
                          fontSize: 22,
                          fontWeight: FontWeight.bold),
                      textAlign: TextAlign.center,
                    ),
                    const SizedBox(height: 8),
                    Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        _buildScoreBadge(
                            'Health', '${result.healthScore}%', Colors.white),
                        const SizedBox(width: 12),
                        _buildScoreBadge(
                            'Confidence',
                            '${(result.confidence * 100).round()}%',
                            Colors.white),
                        const SizedBox(width: 12),
                        _buildScoreBadge('Severity',
                            result.severity.toUpperCase(), Colors.white),
                      ],
                    ),
                  ],
                ),
              ),
              const SizedBox(height: 20),

              // Message
              if (result.message.isNotEmpty) ...[
                const Text('Analysis Summary',
                    style:
                        TextStyle(fontSize: 16, fontWeight: FontWeight.bold)),
                const SizedBox(height: 8),
                Container(
                  width: double.infinity,
                  padding: const EdgeInsets.all(12),
                  decoration: BoxDecoration(
                    color: Colors.grey.shade100,
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: Text(result.message),
                ),
                const SizedBox(height: 20),
              ],

              // Top Predictions
              if (result.topPredictions.isNotEmpty) ...[
                const Text('Top Predictions',
                    style:
                        TextStyle(fontSize: 16, fontWeight: FontWeight.bold)),
                const SizedBox(height: 8),
                ...result.topPredictions.take(5).map((p) => Padding(
                      padding: const EdgeInsets.symmetric(vertical: 4),
                      child: Row(
                        children: [
                          Expanded(child: Text(p.disease)),
                          Container(
                            padding: const EdgeInsets.symmetric(
                                horizontal: 10, vertical: 4),
                            decoration: BoxDecoration(
                              color: Colors.blue.shade50,
                              borderRadius: BorderRadius.circular(12),
                            ),
                            child: Text(
                              '${(p.confidence * 100).round()}%',
                              style: TextStyle(
                                  color: Colors.blue.shade700,
                                  fontWeight: FontWeight.bold,
                                  fontSize: 12),
                            ),
                          ),
                        ],
                      ),
                    )),
                const SizedBox(height: 20),
              ],

              // Save Button
              SizedBox(
                width: double.infinity,
                height: 52,
                child: ElevatedButton(
                  onPressed: () => _saveHealthRecord(result, ctx),
                  style: ElevatedButton.styleFrom(
                      backgroundColor: Theme.of(context).primaryColor,
                      foregroundColor: Colors.white),
                  child: const Text('Save Health Record',
                      style:
                          TextStyle(fontSize: 16, fontWeight: FontWeight.bold)),
                ),
              ),
              const SizedBox(height: 16),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildScoreBadge(String label, String value, Color textColor) {
    return Column(
      children: [
        Text(value,
            style: TextStyle(
                color: textColor, fontSize: 16, fontWeight: FontWeight.bold)),
        Text(label,
            style: TextStyle(color: textColor.withAlpha(200), fontSize: 11)),
      ],
    );
  }

  Future<void> _saveHealthRecord(ScanResult result, BuildContext ctx) async {
    if (widget.cropId == null) {
      AppToast.warning('No crop ID - cannot save record');
      return;
    }

    final recordData = {
      'crop_id': widget.cropId,
      'health_score': result.healthScore,
      'image_url': '', // TODO: Upload image to storage
      'disease_detected': result.isHealthy ? null : result.diseaseName,
      'confidence': result.confidence,
      'severity': result.severity,
      'notes': result.message,
    };

    final saved = await ref
        .read(healthRecordsRepositoryProvider)
        .createHealthRecord(recordData);

    if (saved != null) {
      // If disease detected, create AI-generated treatment recommendation
      if (!result.isHealthy && result.diseaseName.isNotEmpty) {
        await _createAITreatment(result, saved.id);
      }

      // Invalidate providers to force refresh when returning to crop detail
      ref.invalidate(healthRecordsByCropProvider(widget.cropId!));
      ref.invalidate(treatmentsByCropProvider(widget.cropId!));

      AppToast.success('Health record saved!');
      Navigator.of(ctx).pop(); // Close bottom sheet
      if (mounted && context.canPop()) {
        context.pop(); // Go back to crop detail
      }
    } else {
      AppToast.error('Failed to save record');
    }
  }

  Future<void> _createAITreatment(
      ScanResult result, String healthRecordId) async {
    // Generate treatment recommendation based on disease
    final treatmentRecommendation = _getTreatmentForDisease(result.diseaseName);

    final treatmentData = {
      'crop_id': widget.cropId,
      'health_record_id': healthRecordId,
      'disease_name': result.diseaseName,
      'treatment_name': treatmentRecommendation['name'],
      'treatment_type': treatmentRecommendation['type'],
      'application_date': DateTime.now().toIso8601String().split('T')[0],
      'notes': 'AI Recommended: ${treatmentRecommendation['notes']}',
    };

    await ref.read(treatmentsRepositoryProvider).createTreatment(treatmentData);
  }

  Map<String, String> _getTreatmentForDisease(String diseaseName) {
    // AI-based treatment recommendations for common diseases
    final treatments = {
      'Strawberry with Leaf Scorch': {
        'name': 'Copper-based Fungicide',
        'type': 'chemical',
        'notes':
            'Apply copper fungicide weekly. Remove infected leaves. Ensure good air circulation.',
      },
      'Grape with Black Rot': {
        'name': 'Mancozeb Spray',
        'type': 'chemical',
        'notes':
            'Apply protective fungicide before infection. Remove mummified berries. Prune for airflow.',
      },
      'Tomato with Late Blight': {
        'name': 'Chlorothalonil Treatment',
        'type': 'chemical',
        'notes':
            'Apply fungicide at first signs. Remove infected plants. Avoid overhead watering.',
      },
      'Apple with Black Rot': {
        'name': 'Captan Fungicide',
        'type': 'chemical',
        'notes':
            'Spray during growing season. Remove cankers and mummified fruits. Maintain tree health.',
      },
      'Bell Pepper with Bacterial Spot': {
        'name': 'Copper Hydroxide',
        'type': 'chemical',
        'notes':
            'Apply copper bactericide. Use disease-free seeds. Rotate crops yearly.',
      },
      'Squash with Powdery Mildew': {
        'name': 'Neem Oil Spray',
        'type': 'organic',
        'notes':
            'Apply neem oil weekly. Improve air circulation. Water at soil level.',
      },
      'Cedar Apple Rust': {
        'name': 'Myclobutanil Fungicide',
        'type': 'chemical',
        'notes':
            'Apply fungicide in spring. Remove nearby cedar trees if possible.',
      },
    };

    // Default treatment if disease not in list
    return treatments[diseaseName] ??
        {
          'name': 'Consult Agronomist',
          'type': 'organic',
          'notes':
              'Disease detected: $diseaseName. Please consult with a local agronomist for specific treatment recommendations.',
        };
  }

  @override
  Widget build(BuildContext context) {
    // Show error if failed
    if (_errorMessage != null) {
      return Scaffold(
        backgroundColor: Colors.black,
        body: Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              const Icon(Icons.error_outline, color: Colors.white, size: 48),
              const SizedBox(height: 16),
              Text(_errorMessage!,
                  style: const TextStyle(color: Colors.white),
                  textAlign: TextAlign.center),
              const SizedBox(height: 16),
              ElevatedButton(
                onPressed: () {
                  if (context.canPop()) {
                    context.pop();
                  } else {
                    context.go('/crops');
                  }
                },
                child: const Text('Go Back'),
              ),
            ],
          ),
        ),
      );
    }

    // Gallery mode with processing
    if (widget.sourceType == ImageSourceType.gallery) {
      return Scaffold(
        backgroundColor: Colors.black,
        body: Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              if (_isProcessing) ...[
                const CircularProgressIndicator(color: Colors.white),
                const SizedBox(height: 24),
                const Text('Analyzing image...',
                    style: TextStyle(color: Colors.white, fontSize: 18)),
              ] else if (_capturedImage != null) ...[
                ClipRRect(
                  borderRadius: BorderRadius.circular(16),
                  child: Image.file(_capturedImage!,
                      width: 200, height: 200, fit: BoxFit.cover),
                ),
              ] else ...[
                const CircularProgressIndicator(color: Colors.white),
                const SizedBox(height: 16),
                const Text('Opening gallery...',
                    style: TextStyle(color: Colors.white)),
              ]
            ],
          ),
        ),
      );
    }

    // Camera mode
    return Scaffold(
      backgroundColor: Colors.black,
      body: Stack(
        children: [
          // Camera Preview
          if (_initializeControllerFuture != null)
            FutureBuilder<void>(
              future: _initializeControllerFuture,
              builder: (context, snapshot) {
                if (snapshot.connectionState == ConnectionState.done &&
                    _controller != null) {
                  return Center(child: CameraPreview(_controller!));
                } else {
                  return const Center(child: CircularProgressIndicator());
                }
              },
            ),

          // Controls
          Positioned(
            bottom: 0,
            left: 0,
            right: 0,
            child: Container(
              padding: const EdgeInsets.all(24),
              color: Colors.black54,
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceAround,
                children: [
                  IconButton(
                      onPressed: () {
                        if (context.canPop()) {
                          context.pop();
                        } else {
                          context.go('/crops');
                        }
                      },
                      icon: const Icon(Icons.close,
                          color: Colors.white, size: 32)),
                  GestureDetector(
                    onTap: _isProcessing ? null : _takePicture,
                    child: Container(
                      width: 72,
                      height: 72,
                      decoration: BoxDecoration(
                        shape: BoxShape.circle,
                        border: Border.all(color: Colors.white, width: 4),
                        color: Colors.white24,
                      ),
                      child: _isProcessing
                          ? const Center(
                              child: CircularProgressIndicator(
                                  color: Colors.white))
                          : null,
                    ),
                  ),
                  IconButton(
                      onPressed: _pickFromGallery,
                      icon: const Icon(Icons.photo_library,
                          color: Colors.white, size: 32)),
                ],
              ),
            ),
          )
        ],
      ),
    );
  }
}
