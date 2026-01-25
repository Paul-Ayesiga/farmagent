import 'dart:async';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../data/crops_repository.dart';
import '../domain/crop.dart';

// AsyncNotifier to handle list state and operations
class CropsNotifier extends AsyncNotifier<List<Crop>> {
  @override
  FutureOr<List<Crop>> build() {
    return ref.read(cropsRepositoryProvider).getCrops();
  }

  Future<void> refresh() async {
    state = const AsyncLoading();
    state = await AsyncValue.guard(
        () => ref.read(cropsRepositoryProvider).getCrops());
  }

  Future<bool> addCrop(Map<String, dynamic> cropData) async {
    final result = await ref.read(cropsRepositoryProvider).addCrop(cropData);
    if (result != null) {
      // Refresh list on success
      ref.invalidateSelf();
      await future;
      return true;
    }
    return false;
  }
}

final cropsProvider =
    AsyncNotifierProvider<CropsNotifier, List<Crop>>(CropsNotifier.new);
