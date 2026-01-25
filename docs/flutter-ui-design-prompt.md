# FarmAgent Flutter UI Design Prompt

## For AI Design Tools (v0.dev, Claude, ChatGPT, Midjourney, etc.)

---

## Master Prompt for Complete App Design

```
Design a mobile app UI for "FarmAgent" - an AI-powered agricultural assistant 
for Ugandan farmers. The app helps farmers diagnose crop diseases using their 
phone camera, get treatment recommendations, check market prices, and connect 
to buyers.

TARGET USERS:
- Small-scale farmers in Uganda (1-5 acres)
- Age: 25-60 years old
- Device: Low-to-mid range Android phones (majority)
- Network: Often offline or slow 2G/3G
- Tech literacy: Low to medium
- Languages: English and Luganda (local language)

DESIGN REQUIREMENTS:

1. COLOR PALETTE:
   - Primary: Green (#10b981) - represents agriculture, growth, health
   - Secondary: Blue (#3b82f6) - trust, technology
   - Accent: Orange (#f97316) - warnings, attention
   - Success: Bright Green (#22c55e)
   - Warning: Yellow (#eab308)
   - Danger: Red (#ef4444)
   - Background: Light gray (#f3f4f6)
   - Text: Dark gray (#1f2937)

2. TYPOGRAPHY:
   - Headers: Bold, large (24-32px)
   - Body: Regular, readable (16-18px) - important for older farmers
   - Buttons: Medium weight, clear (16-18px)
   - Font: Use system fonts (Roboto for Android) for better performance

3. ICONOGRAPHY:
   - Use simple, recognizable icons
   - Agriculture-related: 🌾🌽🌿☕🥔🫘 (or similar clean SVG versions)
   - Function icons: Camera, Market, Wallet, Bell, User
   - Status icons: Health bars, checkmarks, warnings

4. LAYOUT PRINCIPLES:
   - Large touch targets (minimum 48x48dp) - for outdoor use with gloves
   - High contrast for outdoor visibility
   - Bottom navigation for main features (thumb-friendly)
   - Generous spacing and padding
   - Offline-first indicators

5. KEY SCREENS TO DESIGN:

   A. SPLASH/ONBOARDING SCREEN
      - FarmAgent logo
      - Tagline: "Your AI Farm Assistant"
      - Language selection (English/Luganda)
      - Simple 3-slide tutorial

   B. HOME DASHBOARD
      - Welcome message with farmer's name
      - Weather alert card (with icon and brief message)
      - "My Crops" section:
        * List of farmer's crops with health scores
        * Visual health indicators (progress bars/emoji)
        * Quick status: Healthy / Needs Attention / Critical
      - Quick action buttons:
        * Large "Scan Crop" button (camera icon)
        * "Market Prices" button
        * "My Orders" button
      - Recent notifications preview
      - Bottom navigation: Home, Scan, Market, Messages, Profile

   C. CROP SCANNING SCREEN
      - Full-screen camera view
      - Overlay guide (frame showing where to position crop)
      - "Take Photo" button (large, centered bottom)
      - Gallery upload option
      - Tips overlay: "Focus on affected leaves"
      - Loading state: "Analyzing your crop..." with animation

   D. DIAGNOSIS RESULT SCREEN
      - Large crop image at top
      - Disease name card:
        * Disease name in bold
        * Confidence score (87%)
        * Severity badge (Mild/Moderate/Severe)
      - Affected area visualization
      - Treatment recommendation card:
        * Primary treatment option
        * Cost estimate
        * Where to buy
        * How to apply
      - Prevention tips section
      - Action buttons:
        * "Order Treatment" (primary CTA)
        * "Save to History"
        * "Set Reminder"
        * "Consult Expert"

   E. MARKET PRICES SCREEN
      - Header: "Current Market Prices" with last updated time
      - Filter chips: All / My Crops / Top Prices
      - Price cards for each commodity:
        * Crop icon and name
        * Current price per kg
        * Price trend (up/down with percentage)
        * Graph showing 7-day trend
        * "Sell Now" button
      - Market insights card at bottom

   F. SELL REQUEST SCREEN
      - Crop selection dropdown
      - Quantity input (kg)
      - Quality selector (Premium/Standard/Mixed)
      - Expected harvest date picker
      - Location (auto-filled from GPS)
      - Upload crop photos
      - Price expectation (with market suggestion)
      - "Find Buyers" button

   G. MESSAGES/NOTIFICATIONS SCREEN
      - Tab switcher: All / Alerts / Recommendations
      - Notification cards:
        * Icon (based on type)
        * Title and message
        * Time stamp
        * Action button if applicable
      - Mark all as read option

   H. PROFILE SCREEN
      - Farmer profile header:
        * Profile photo
        * Name
        * Farm size
        * Location
      - Stats cards:
        * Total crops
        * Average health score
        * Successful harvests
      - Settings sections:
        * Language preference
        * Notification settings
        * Payment methods
        * Help & Support
      - Logout button

6. SPECIAL FEATURES TO SHOW:

   - OFFLINE INDICATOR:
     * Small banner when offline
     * "Syncing..." status when back online
   
   - HEALTH SCORE VISUALIZATION:
     * Use progress bars with colors:
       - 80-100%: Green (Healthy)
       - 50-79%: Yellow (Fair)
       - 0-49%: Red (Poor)
   
   - EMPTY STATES:
     * "No crops added yet" with illustration
     * "Add your first crop" CTA button
   
   - LOADING STATES:
     * Skeleton screens for content loading
     * Shimmer effect on cards
     * "Analyzing image..." with progress indicator
   
   - ERROR STATES:
     * Friendly illustrations
     * Clear error message
     * "Try Again" button

7. ACCESSIBILITY:
   - High contrast ratios (WCAG AA minimum)
   - Large text options
   - Screen reader friendly labels
   - Color-blind safe palette
   - Works in bright sunlight

8. LOCALIZATION:
   - Show both English and Luganda text examples
   - Language toggle in header/profile
   - Date formats: DD/MM/YYYY
   - Currency: UGX (Uganda Shillings)

9. ANIMATIONS (SUBTLE):
   - Page transitions: Slide from right
   - Card interactions: Gentle lift on tap
   - Success states: Checkmark animation
   - Loading: Gentle pulse or shimmer
   - Keep animations short (200-300ms)

10. BRANDING:
    - App name: "FarmAgent"
    - Icon: Stylized leaf with AI circuit pattern
    - Tagline: "Your AI Farm Assistant"
    - Voice: Friendly, helpful, empowering

DELIVERABLES REQUESTED:
- High-fidelity mockups for all 8 key screens
- Light mode design (dark mode optional)
- Mobile viewport: 360x800px (standard Android)
- Component library showing: buttons, cards, inputs, icons
- Color palette with hex codes
- Typography scale
- Spacing system
- Interactive prototype (if possible)

DESIGN STYLE:
- Modern but not overly minimalist
- Approachable and friendly
- Trust-inspiring
- Clear information hierarchy
- Action-oriented (clear CTAs)
- Visual but not cluttered

INSPIRATION REFERENCES:
- WhatsApp: Familiar, simple chat interface
- M-Pesa/Mobile Money apps: Trust, clarity, transaction flow
- Weather apps: Clear data visualization, icons
- Health apps: Progress tracking, scores, recommendations

Please create a comprehensive UI design that makes crop disease diagnosis 
and farm management accessible to Ugandan farmers with varying levels of 
technology literacy.
```

---

## Screen-Specific Prompts

### If you want to design screens individually:

#### 1. Home Dashboard Prompt

```
Design a home dashboard screen for FarmAgent, an agricultural AI app for 
Ugandan farmers.

SCREEN ELEMENTS:
- Top header:
  * "Welcome back, [Farmer Name]" greeting
  * Weather icon and temperature
  * Notification bell icon (with badge if unread)
  * Language toggle (EN/LG)

- Weather alert card:
  * Cloud/sun/rain icon
  * "Rain expected in 2 days - Good time to fertilize!"
  * Light blue background
  * Dismissible

- "My Crops" section:
  * Section header with "View All" link
  * 2-3 crop cards showing:
    - Crop emoji/icon (🌽 maize, 🌿 cassava, ☕ coffee)
    - Crop name and field name
    - Health score: 85% with green progress bar
    - Last checked: "2 days ago"
    - Status badge: "Healthy" (green) or "Needs Attention" (orange)
  * Add crop button (+)

- Quick Actions (2 large buttons):
  * "📸 Scan Crop" - Primary green button
  * "📈 Market Prices" - Secondary blue button

- Recent Activity feed:
  * 2-3 notification previews
  * "View All" link

- Bottom Navigation Bar:
  * Home (active - green)
  * Scan (camera icon)
  * Market (chart icon)
  * Messages (chat icon)
  * Profile (user icon)

DESIGN SPECS:
- Viewport: 360x800px
- Colors: Green primary #10b981, Blue secondary #3b82f6
- Spacing: Generous (16-24px between sections)
- Cards: White with subtle shadow, 8px border radius
- Typography: Bold headers 24px, body 16px

Make it feel welcoming, easy to scan, and action-oriented.
```

#### 2. Camera/Scan Screen Prompt

```
Design a crop scanning/camera screen for FarmAgent.

SCREEN ELEMENTS:
- Full-screen camera viewfinder
- Top bar (semi-transparent):
  * Back button (left)
  * Flash toggle (right)
  * Gallery icon (right)

- Center overlay:
  * Rectangular frame guide (dashed border)
  * Text: "Position crop leaves within frame"
  * Icon showing example positioning

- Bottom section:
  * Large circular "Capture" button (white, centered)
  * "Upload from Gallery" text button below

- Tips panel (collapsible):
  * 📸 "Take in good lighting"
  * 🌿 "Focus on affected leaves"
  * 📏 "Get close but not blurry"

STATES TO SHOW:
1. Camera ready state
2. Capturing (button animation)
3. Processing: "Analyzing your crop..." with spinner
4. Success: Brief checkmark animation before results

DESIGN SPECS:
- Full screen, dark semi-transparent overlays
- Capture button: 80x80px, white with subtle shadow
- Guide frame: Green dashed border, 280x280px
- Text: White for visibility over camera
- Icons: Simple, outlined style

Make it intuitive for farmers who may not use camera apps often.
```

#### 3. Diagnosis Results Prompt

```
Design a disease diagnosis results screen for FarmAgent.

SCREEN LAYOUT:
- Header:
  * Back button
  * "Diagnosis Results" title
  * Share icon

- Crop Image:
  * Full-width image of the crop photo taken
  * Rounded corners
  * Affected areas highlighted with red overlay

- Diagnosis Card (prominent):
  * Disease icon (warning triangle)
  * Disease name: "Cassava Brown Streak Disease"
  * Confidence: 87% (with circular progress indicator)
  * Severity badge: "Moderate" (orange pill)
  * Affected area: 35% of leaves

- Treatment Section:
  * "💊 Recommended Treatment" header
  * Treatment card:
    - Treatment name: "Neem-based organic pesticide"
    - Application: "Spray every 7 days for 3 weeks"
    - Dosage: "50ml per 10 liters of water"
    - Estimated cost: UGX 15,000
    - Availability: "Kampala Farmers Market"
  * "Order Treatment" button (primary green)

- Prevention Tips (collapsible):
  * "🛡️ How to Prevent This"
  * Bulleted tips
  * Light blue background

- Action Buttons:
  * "Set 7-Day Reminder" (outline button)
  * "Consult Expert" (outline button)
  * "Save to History" (text button)

DESIGN SPECS:
- Scrollable content
- Cards: White background, 12px border radius, subtle shadow
- Orange for warnings/moderate severity
- Green for action buttons
- Icons: Simple, meaningful
- Padding: 16px around cards

Make the treatment information crystal clear and actionable.
```

#### 4. Market Prices Prompt

```
Design a market prices screen for FarmAgent.

SCREEN LAYOUT:
- Header:
  * "Market Prices" title
  * Filter icon
  * Last updated: "2 hours ago"

- Filter Chips (horizontal scroll):
  * All (active)
  * My Crops
  * Trending
  * Near Me

- Price Cards (scrollable list):
  Each card shows:
  * Crop icon (🌽🌿☕)
  * Crop name: "Maize"
  * Current price: "UGX 1,200 per kg"
  * Trend indicator:
    - Arrow up/down
    - Percentage: "+5%" (green) or "-3%" (red)
  * Mini sparkline chart (7-day trend)
  * Location: "Kampala Market"
  * "Sell Now" button (outline)

- Market Insights Card (at bottom):
  * "💡 Market Insights"
  * "Beans prices up 8% this week. Good time to sell!"
  * Light yellow background

- Floating Action Button:
  * "+" button (bottom right)
  * "Create Sell Request"

DESIGN SPECS:
- Clean, scannable list
- Cards: White, 8px radius, 12px spacing
- Trend colors: Green for up, Red for down
- Chart: Simple line, 60px tall, gray with accent color
- Price: Large, bold (20px)
- CTAs: Clear, action-oriented

Make price comparison quick and trends immediately visible.
```

#### 5. Profile Screen Prompt

```
Design a profile/settings screen for FarmAgent.

SCREEN LAYOUT:
- Profile Header:
  * Profile photo (circular, 80x80px)
  * Farmer name: "John Mukasa"
  * Edit icon
  * Location: "Wakiso District"
  * Member since: "Jan 2026"

- Stats Row (3 cards):
  * "🌾 5 Crops" - Active crops
  * "📊 82% Avg Health" - Overall health
  * "✅ 12 Harvests" - Completed

- Settings Sections:

  "Account Settings"
  * My Farm Details
  * Payment Methods
  * Language: English ▼

  "Preferences"
  * Notifications (toggle switches):
    - Disease alerts
    - Market updates  
    - Harvest reminders
    - Expert messages

  "Support"
  * Help Center
  * Contact Extension Officer
  * Tutorial Videos
  * Share App

  "About"
  * App Version: 1.0.0
  * Terms of Service
  * Privacy Policy

- Logout Button (bottom):
  * Red text
  * Confirmation required

DESIGN SPECS:
- Profile header: Green gradient background
- Stats cards: White, icons with accent colors
- Settings: Grouped lists, dividers between sections
- List items: Left icon, title, right arrow
- Toggles: Standard material switches
- Logout: Destructive red color

Make settings organized and easy to find.
```

---

## Component Library Prompt

```
Create a design system/component library for FarmAgent with these components:

BUTTONS:
1. Primary Button:
   - Background: #10b981 (green)
   - Text: White, 16px medium
   - Height: 48px
   - Border radius: 8px
   - States: Default, Pressed, Disabled

2. Secondary Button:
   - Border: 2px solid #10b981
   - Text: #10b981, 16px medium
   - Background: White
   - Same sizing as primary

3. Text Button:
   - No background
   - Text: #10b981, 16px medium
   - Underline on hover

CARDS:
1. Standard Card:
   - Background: White
   - Border radius: 12px
   - Shadow: 0 2px 8px rgba(0,0,0,0.08)
   - Padding: 16px

2. Alert Card:
   - Left border: 4px colored (blue/yellow/red/green)
   - Icon on left
   - Close button on right

3. Crop Card:
   - Horizontal layout
   - Crop icon/image left (48x48px)
   - Info middle (name, status)
   - Health score right
   - Progress bar bottom

INPUT FIELDS:
- Height: 48px
- Border: 1px solid #d1d5db
- Border radius: 8px
- Padding: 12px 16px
- Focus: Border #10b981
- Label: 14px, #6b7280, above input
- Error: Red border + error message below

ICONS:
- Style: Outlined
- Size: 24x24px (default), 20x20px (small), 32x32px (large)
- Color: Inherit from context
- Agriculture: Leaf, plant, tractor, money, weather
- UI: Home, camera, chart, message, user, settings, arrow

BADGES:
- Pills: Border radius 12px
- Colors based on status:
  * Success: Green background, dark green text
  * Warning: Yellow background, dark yellow text
  * Danger: Red background, white text
  * Info: Blue background, dark blue text
- Padding: 4px 12px
- Font: 12px medium

SPACING SYSTEM:
- 4px, 8px, 12px, 16px, 24px, 32px, 48px
- Use 16px as base unit

ELEVATION (Shadows):
- Level 1: 0 1px 3px rgba(0,0,0,0.1)
- Level 2: 0 2px 8px rgba(0,0,0,0.08)
- Level 3: 0 4px 16px rgba(0,0,0,0.12)

TYPOGRAPHY SCALE:
- H1: 32px bold
- H2: 24px bold
- H3: 20px semibold
- Body: 16px regular
- Caption: 14px regular
- Small: 12px regular

Create mockups showing all these components in use.
```

---

## Quick Figma/Design Tool Prompts

### For Figma AI/Plugins:

```
Create FarmAgent app screens in Figma:
- 360x800px frames
- Android Material Design
- Primary color: #10b981 green
- Include: Home, Camera, Results, Market, Profile screens
- Components: Bottom nav, cards, buttons, inputs
- Auto-layout for responsive design
```

### For v0.dev:

```
Build a Flutter/React Native design system for FarmAgent agricultural app.
Theme: Green (#10b981), professional, farmer-friendly.
Screens needed: Dashboard with crop health cards, camera scanner, 
disease results, market prices, profile.
Include: Large buttons, progress indicators, status badges, navigation bar.
Target: Android, offline-capable, high contrast.
```

### For Midjourney (Visual Inspiration):

```
Mobile app UI design for agricultural assistant, clean modern interface,
green and blue color scheme, crop disease detection screens, farmer using
smartphone in field, African context, professional product design, 
Dribbble style, --ar 9:16
```

---

## Usage Tips

1. **Start with Master Prompt** - Get complete app design
2. **Iterate with Screen-Specific** - Refine individual screens
3. **Use Component Prompt** - Build design system
4. **Generate Multiple Variations** - Compare approaches
5. **Combine Tools**:
   - Use AI for quick mockups
   - Refine in Figma/Adobe XD
   - Use Midjourney for visual inspiration
   - Use v0.dev for interactive prototypes

---

## Export Specifications

When getting designs, request:
- PNG exports at 2x resolution (720x1600px)
- SVG for icons and illustrations
- Figma/XD file for developer handoff
- Style guide PDF
- Interactive prototype link
- Design tokens (colors, spacing, typography in JSON)

---

## Next Steps After Design

1. Review designs with potential users (farmers)
2. Test contrast in bright sunlight
3. Validate with low-tech literacy users
4. Ensure all text is readable
5. Check touch target sizes
6. Implement in Flutter using the designs
7. Iterate based on user testing

---

Copy any of these prompts into your preferred AI design tool and adjust as needed!