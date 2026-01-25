import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

class ScaffoldWithNavBar extends StatelessWidget {
  final GoRouterState state;
  final Widget child;

  const ScaffoldWithNavBar({
    super.key,
    required this.state,
    required this.child,
  });

  @override
  Widget build(BuildContext context) {
    // Determine selected index based on path
    int selectedIndex = 0;
    final path = state.uri.path;

    if (path.startsWith('/fields')) {
      selectedIndex = 1;
    } else if (path.startsWith('/crops')) {
      selectedIndex = 2;
    } else if (path.startsWith('/payments')) {
      selectedIndex = 3;
    } else if (path.startsWith('/dashboard') || path == '/') {
      selectedIndex = 0;
    }

    final primaryGreen = Theme.of(context).primaryColor;

    return Scaffold(
      body: child,
      extendBody: true, // Required for the FAB notch effect
      floatingActionButton: FloatingActionButton(
        heroTag: 'main_ai_fab',
        onPressed: () => context.push('/ai/chat'),
        backgroundColor: primaryGreen,
        elevation: 4,
        shape: const CircleBorder(),
        tooltip: 'AI Assistant',
        child: const Icon(Icons.auto_awesome, color: Colors.white, size: 28),
      ),
      floatingActionButtonLocation: FloatingActionButtonLocation.centerDocked,
      bottomNavigationBar: BottomAppBar(
        shape: const CircularNotchedRectangle(),
        notchMargin: 8,
        color: Colors.white,
        elevation: 8,
        height: 56,
        padding: EdgeInsets.zero,
        child: Row(
          mainAxisAlignment: MainAxisAlignment.spaceAround,
          children: [
            // Left side tabs
            _buildNavItem(context, Icons.home_outlined, Icons.home, 'Home', 0,
                selectedIndex, () => context.go('/dashboard')),
            _buildNavItem(context, Icons.map_outlined, Icons.map, 'Fields', 1,
                selectedIndex, () => context.go('/fields')),

            // Spacer for the FAB
            const SizedBox(width: 56),

            // Right side tabs
            _buildNavItem(context, Icons.grass_outlined, Icons.grass, 'Crops',
                2, selectedIndex, () => context.go('/crops')),
            _buildNavItem(context, Icons.storefront_outlined, Icons.storefront,
                'Market', 3, selectedIndex, () => context.go('/payments')),
          ],
        ),
      ),
    );
  }

  Widget _buildNavItem(
      BuildContext context,
      IconData icon,
      IconData selectedIcon,
      String label,
      int index,
      int selectedIndex,
      VoidCallback onTap) {
    final isSelected = selectedIndex == index;
    final primaryGreen = Theme.of(context).primaryColor;

    return Expanded(
      child: InkWell(
        onTap: onTap,
        child: SizedBox(
          height: 56,
          child: Column(
            mainAxisSize: MainAxisSize.min,
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Icon(
                isSelected ? selectedIcon : icon,
                color: isSelected ? primaryGreen : Colors.grey.shade500,
                size: 20,
              ),
              Text(
                label,
                style: TextStyle(
                  fontSize: 9,
                  fontWeight: isSelected ? FontWeight.bold : FontWeight.normal,
                  color: isSelected ? primaryGreen : Colors.grey.shade500,
                ),
                maxLines: 1,
                overflow: TextOverflow.ellipsis,
              ),
            ],
          ),
        ),
      ),
    );
  }
}
