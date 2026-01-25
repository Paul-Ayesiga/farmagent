import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../domain/crop.dart';
import '../domain/health_record.dart';
import '../domain/treatment.dart';
import '../data/crops_repository.dart';
import 'crops_providers.dart';

// Global route observer for detecting navigation changes
final RouteObserver<PageRoute> routeObserver = RouteObserver<PageRoute>();

class CropDetailScreen extends ConsumerStatefulWidget {
  final String cropId;

  const CropDetailScreen({super.key, required this.cropId});

  @override
  ConsumerState<CropDetailScreen> createState() => _CropDetailScreenState();
}

class _CropDetailScreenState extends ConsumerState<CropDetailScreen>
    with SingleTickerProviderStateMixin, RouteAware {
  late TabController _tabController;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 2, vsync: this);
  }

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    // Subscribe to route changes
    final route = ModalRoute.of(context);
    if (route is PageRoute) {
      routeObserver.subscribe(this, route);
    }
  }

  @override
  void dispose() {
    routeObserver.unsubscribe(this);
    _tabController.dispose();
    super.dispose();
  }

  @override
  void didPopNext() {
    // Called when returning to this screen from another route
    // Refresh the data
    ref.invalidate(healthRecordsByCropProvider(widget.cropId));
    ref.invalidate(treatmentsByCropProvider(widget.cropId));
  }

  void _showImageSourcePicker(
      BuildContext context, String cropId, String cropName) {
    showModalBottomSheet(
      context: context,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (ctx) => SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(24),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const Text('Scan Crop Health',
                  style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold)),
              const SizedBox(height: 8),
              const Text('Choose how you want to capture the image',
                  style: TextStyle(color: Colors.grey)),
              const SizedBox(height: 24),
              ListTile(
                leading: Container(
                  padding: const EdgeInsets.all(12),
                  decoration: BoxDecoration(
                    color: Colors.blue.shade50,
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: Icon(Icons.camera_alt, color: Colors.blue.shade600),
                ),
                title: const Text('Take Photo',
                    style: TextStyle(fontWeight: FontWeight.bold)),
                subtitle: const Text('Use camera to capture crop image'),
                trailing: const Icon(Icons.chevron_right),
                onTap: () {
                  Navigator.pop(ctx);
                  context.push(
                      '/crops/$cropId/scan?source=camera&name=${Uri.encodeComponent(cropName)}');
                },
              ),
              const SizedBox(height: 8),
              ListTile(
                leading: Container(
                  padding: const EdgeInsets.all(12),
                  decoration: BoxDecoration(
                    color: Colors.green.shade50,
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child:
                      Icon(Icons.photo_library, color: Colors.green.shade600),
                ),
                title: const Text('Choose from Gallery',
                    style: TextStyle(fontWeight: FontWeight.bold)),
                subtitle: const Text('Upload existing image'),
                trailing: const Icon(Icons.chevron_right),
                onTap: () {
                  Navigator.pop(ctx);
                  context.push(
                      '/crops/$cropId/scan?source=gallery&name=${Uri.encodeComponent(cropName)}');
                },
              ),
              const SizedBox(height: 16),
            ],
          ),
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final cropsAsync = ref.watch(cropsProvider);
    final primaryGreen = Theme.of(context).primaryColor;

    return cropsAsync.when(
      data: (crops) {
        final crop = crops.firstWhere((c) => c.id == widget.cropId,
            orElse: () => Crop(
                id: '',
                name: 'Not Found',
                variety: '',
                fieldId: '',
                plantingDate: '',
                expectedHarvest: '',
                healthScore: 0,
                status: '',
                imageUrl: ''));

        if (crop.id.isEmpty) {
          return Scaffold(
            appBar: AppBar(title: const Text('Crop Details')),
            body: const Center(child: Text('Crop not found')),
          );
        }

        return Scaffold(
          appBar: AppBar(
            title: Text(crop.name),
            actions: [
              IconButton(
                icon: const Icon(Icons.refresh),
                onPressed: () {
                  ref.invalidate(cropsProvider);
                  ref.invalidate(healthRecordsByCropProvider(widget.cropId));
                  ref.invalidate(treatmentsByCropProvider(widget.cropId));
                },
              ),
            ],
          ),
          body: NestedScrollView(
            headerSliverBuilder: (context, innerBoxIsScrolled) => [
              SliverToBoxAdapter(
                child: Padding(
                  padding: const EdgeInsets.all(16),
                  child: Column(
                    children: [
                      // Header Card
                      Container(
                        padding: const EdgeInsets.all(20),
                        decoration: BoxDecoration(
                          color: Colors.white,
                          borderRadius: BorderRadius.circular(16),
                          boxShadow: const [
                            BoxShadow(color: Colors.black12, blurRadius: 10)
                          ],
                        ),
                        child: Row(
                          children: [
                            Container(
                              padding: const EdgeInsets.all(16),
                              decoration: BoxDecoration(
                                  color: primaryGreen.withAlpha(25),
                                  shape: BoxShape.circle),
                              child: Icon(Icons.eco,
                                  size: 36, color: primaryGreen),
                            ),
                            const SizedBox(width: 16),
                            Expanded(
                              child: Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  Text(crop.name,
                                      style: const TextStyle(
                                          fontSize: 22,
                                          fontWeight: FontWeight.bold)),
                                  Text('${crop.variety} • ${crop.status}',
                                      style: TextStyle(
                                          color: Colors.grey.shade600)),
                                  const SizedBox(height: 8),
                                  Row(
                                    children: [
                                      _buildMiniStat('Health',
                                          '${crop.healthScore}%', primaryGreen),
                                      const SizedBox(width: 16),
                                      _buildMiniStat('Harvest',
                                          crop.expectedHarvest, Colors.orange),
                                    ],
                                  ),
                                ],
                              ),
                            ),
                          ],
                        ),
                      ),
                      const SizedBox(height: 16),

                      // Scan Button
                      SizedBox(
                        width: double.infinity,
                        height: 52,
                        child: ElevatedButton.icon(
                          onPressed: () => _showImageSourcePicker(
                              context, widget.cropId, crop.name),
                          icon: const Icon(Icons.qr_code_scanner),
                          label: const Text('Scan Health'),
                          style: ElevatedButton.styleFrom(
                            backgroundColor: primaryGreen,
                            foregroundColor: Colors.white,
                            textStyle: const TextStyle(
                                fontSize: 16, fontWeight: FontWeight.bold),
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
              ),
              SliverPersistentHeader(
                pinned: true,
                delegate: _SliverTabBarDelegate(
                  TabBar(
                    controller: _tabController,
                    labelColor: primaryGreen,
                    unselectedLabelColor: Colors.grey,
                    indicatorColor: primaryGreen,
                    tabs: const [
                      Tab(
                          icon: Icon(Icons.medical_information),
                          text: 'Health Records'),
                      Tab(icon: Icon(Icons.healing), text: 'Treatments'),
                    ],
                  ),
                ),
              ),
            ],
            body: TabBarView(
              controller: _tabController,
              children: [
                // Health Records Tab
                _HealthRecordsTab(cropId: widget.cropId),
                // Treatments Tab
                _TreatmentsTab(cropId: widget.cropId),
              ],
            ),
          ),
        );
      },
      loading: () => Scaffold(
        appBar: AppBar(title: const Text('Crop Details')),
        body: const Center(child: CircularProgressIndicator()),
      ),
      error: (e, s) => Scaffold(
        appBar: AppBar(title: const Text('Crop Details')),
        body: Center(child: Text('Error: $e')),
      ),
    );
  }

  Widget _buildMiniStat(String label, String value, Color color) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(label,
            style: TextStyle(fontSize: 11, color: Colors.grey.shade500)),
        Text(value,
            style: TextStyle(
                fontSize: 14, fontWeight: FontWeight.bold, color: color)),
      ],
    );
  }
}

// Sliver delegate for sticky TabBar
class _SliverTabBarDelegate extends SliverPersistentHeaderDelegate {
  final TabBar tabBar;

  _SliverTabBarDelegate(this.tabBar);

  @override
  double get minExtent => tabBar.preferredSize.height;
  @override
  double get maxExtent => tabBar.preferredSize.height;

  @override
  Widget build(
      BuildContext context, double shrinkOffset, bool overlapsContent) {
    return Container(
      color: Colors.white,
      child: tabBar,
    );
  }

  @override
  bool shouldRebuild(_SliverTabBarDelegate oldDelegate) => false;
}

// Providers for health records and treatments
final healthRecordsByCropProvider =
    FutureProvider.family<List<HealthRecord>, String>((ref, cropId) async {
  final repo = ref.read(cropsRepositoryProvider);
  return repo.getHealthRecords(cropId);
});

final treatmentsByCropProvider =
    FutureProvider.family<List<Treatment>, String>((ref, cropId) async {
  final repo = ref.read(cropsRepositoryProvider);
  return repo.getTreatments(cropId);
});

// ==================== Health Records Tab ====================
class _HealthRecordsTab extends ConsumerWidget {
  final String cropId;

  const _HealthRecordsTab({required this.cropId});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final recordsAsync = ref.watch(healthRecordsByCropProvider(cropId));

    return recordsAsync.when(
      data: (records) {
        if (records.isEmpty) {
          return _buildEmptyState(
            icon: Icons.medical_information_outlined,
            title: 'No health records yet',
            subtitle: 'Scan your crop to create a health record',
          );
        }

        return ListView.builder(
          padding: const EdgeInsets.all(16),
          itemCount: records.length,
          itemBuilder: (context, index) => _buildRecordCard(records[index]),
        );
      },
      loading: () => const Center(child: CircularProgressIndicator()),
      error: (e, s) => Center(
          child: Text('Error: $e', style: const TextStyle(color: Colors.red))),
    );
  }

  Widget _buildRecordCard(HealthRecord record) {
    final isHealthy = record.isHealthy;

    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      elevation: 0,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(12),
        side: BorderSide(color: Colors.grey.shade200),
      ),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Row(
          children: [
            Container(
              padding: const EdgeInsets.all(12),
              decoration: BoxDecoration(
                color: isHealthy ? Colors.green.shade50 : Colors.orange.shade50,
                borderRadius: BorderRadius.circular(12),
              ),
              child: Icon(
                isHealthy ? Icons.check_circle : Icons.warning,
                color: isHealthy ? Colors.green : Colors.orange,
              ),
            ),
            const SizedBox(width: 16),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    isHealthy
                        ? 'Healthy'
                        : (record.diseaseDetected ?? 'Unknown Issue'),
                    style: const TextStyle(
                        fontWeight: FontWeight.bold, fontSize: 15),
                  ),
                  const SizedBox(height: 4),
                  Text(
                    'Health Score: ${record.healthScore}%',
                    style: TextStyle(color: Colors.grey.shade600, fontSize: 12),
                  ),
                ],
              ),
            ),
            Column(
              crossAxisAlignment: CrossAxisAlignment.end,
              children: [
                Text(
                  _formatDate(record.createdAt),
                  style: TextStyle(color: Colors.grey.shade500, fontSize: 11),
                ),
                if (record.confidence != null)
                  Container(
                    margin: const EdgeInsets.only(top: 4),
                    padding:
                        const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                    decoration: BoxDecoration(
                      color: Colors.blue.shade50,
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: Text(
                      '${(record.confidence! * 100).round()}%',
                      style: TextStyle(
                          color: Colors.blue.shade700,
                          fontSize: 10,
                          fontWeight: FontWeight.bold),
                    ),
                  ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildEmptyState(
      {required IconData icon,
      required String title,
      required String subtitle}) {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(32),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(icon, size: 64, color: Colors.grey.shade300),
            const SizedBox(height: 16),
            Text(title,
                style: TextStyle(
                    fontSize: 16,
                    color: Colors.grey.shade600,
                    fontWeight: FontWeight.w500)),
            const SizedBox(height: 4),
            Text(subtitle,
                style: TextStyle(fontSize: 13, color: Colors.grey.shade400)),
          ],
        ),
      ),
    );
  }

  String _formatDate(DateTime date) {
    return '${date.day}/${date.month}/${date.year}';
  }
}

// ==================== Treatments Tab ====================
class _TreatmentsTab extends ConsumerWidget {
  final String cropId;

  const _TreatmentsTab({required this.cropId});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final treatmentsAsync = ref.watch(treatmentsByCropProvider(cropId));

    return treatmentsAsync.when(
      data: (treatments) {
        if (treatments.isEmpty) {
          return _buildEmptyState(
            icon: Icons.healing_outlined,
            title: 'No treatments yet',
            subtitle: 'Treatments will appear here after AI analysis',
          );
        }

        return ListView.builder(
          padding: const EdgeInsets.all(16),
          itemCount: treatments.length,
          itemBuilder: (context, index) =>
              _buildTreatmentCard(treatments[index]),
        );
      },
      loading: () => const Center(child: CircularProgressIndicator()),
      error: (e, s) => Center(
          child: Text('Error: $e', style: const TextStyle(color: Colors.red))),
    );
  }

  Widget _buildTreatmentCard(Treatment treatment) {
    final isOrganic = treatment.isOrganic;

    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      elevation: 0,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(12),
        side: BorderSide(color: Colors.grey.shade200),
      ),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Container(
                  padding: const EdgeInsets.all(10),
                  decoration: BoxDecoration(
                    color: isOrganic
                        ? Colors.green.shade50
                        : Colors.purple.shade50,
                    borderRadius: BorderRadius.circular(10),
                  ),
                  child: Icon(
                    isOrganic ? Icons.eco : Icons.science,
                    color: isOrganic ? Colors.green : Colors.purple,
                    size: 20,
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(treatment.treatmentName,
                          style: const TextStyle(
                              fontWeight: FontWeight.bold, fontSize: 15)),
                      Text('For: ${treatment.diseaseName}',
                          style: TextStyle(
                              color: Colors.grey.shade600, fontSize: 12)),
                    ],
                  ),
                ),
                Container(
                  padding:
                      const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                  decoration: BoxDecoration(
                    color: isOrganic
                        ? Colors.green.shade100
                        : Colors.purple.shade100,
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: Text(
                    isOrganic ? 'Organic' : 'Chemical',
                    style: TextStyle(
                      color: isOrganic
                          ? Colors.green.shade800
                          : Colors.purple.shade800,
                      fontSize: 10,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                ),
              ],
            ),
            if (treatment.notes != null && treatment.notes!.isNotEmpty) ...[
              const SizedBox(height: 12),
              Text(treatment.notes!,
                  style: TextStyle(color: Colors.grey.shade700, fontSize: 13)),
            ],
            const SizedBox(height: 12),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(
                  'Applied: ${_formatDate(treatment.applicationDate)}',
                  style: TextStyle(color: Colors.grey.shade500, fontSize: 11),
                ),
                if (treatment.effectiveness != null)
                  Row(
                    children: List.generate(
                        5,
                        (i) => Icon(
                              i < treatment.effectiveness!
                                  ? Icons.star
                                  : Icons.star_border,
                              size: 14,
                              color: Colors.amber,
                            )),
                  ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildEmptyState(
      {required IconData icon,
      required String title,
      required String subtitle}) {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(32),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(icon, size: 64, color: Colors.grey.shade300),
            const SizedBox(height: 16),
            Text(title,
                style: TextStyle(
                    fontSize: 16,
                    color: Colors.grey.shade600,
                    fontWeight: FontWeight.w500)),
            const SizedBox(height: 4),
            Text(subtitle,
                style: TextStyle(fontSize: 13, color: Colors.grey.shade400)),
          ],
        ),
      ),
    );
  }

  String _formatDate(DateTime date) {
    return '${date.day}/${date.month}/${date.year}';
  }
}
