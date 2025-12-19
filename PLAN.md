# Hatchbridge.com Restyling Plan

## Analysis of Target Design (Hatchbridge.com)

From the website screenshot, the key design elements are:

### Color Palette
- **Primary Brand Color**: Golden Yellow `#FFC421` (matches the logo SVG)
- **Dark Background**: Deep black/charcoal for headers and hero sections
- **Text on Dark**: Pure white `#FFFFFF`
- **Button Style**: Yellow/gold background with dark text
- **Overall Feel**: Modern, professional startup incubator aesthetic

### Typography
- Clean sans-serif fonts
- Bold, large headings
- Good contrast and readability

### UI Elements
- Yellow CTA buttons ("Become a Member", "Join a Program")
- Dark navigation header with white text
- Clean, minimal design with strong branding

---

## Current State Analysis

The current "hatchbridge" theme has:
- Light background (white)
- Dark slate as primary (`220 45% 22%`)
- Gold as secondary (`42 75% 52%`)

**Problem**: This is inverted from the actual website. The real Hatchbridge.com uses:
- Dark backgrounds
- Gold/yellow as the primary accent and CTA color
- Much more prominent use of the brand yellow

---

## Implementation Plan

### Phase 1: Update Hatchbridge Theme Colors

**File**: `frontend/assets/css/main.css`

Update the `.theme-hatchbridge` to match the website:

```css
:root,
.homebox,
.theme-hatchbridge {
  /* Hatchbridge brand - Dark mode with gold accents */
  --background: 220 15% 13%;           /* Dark charcoal #1e2128 */
  --background-accent: 220 18% 10%;    /* Darker accent #181b22 */
  --foreground: 0 0% 98%;              /* Near white */

  --primary: 45 100% 50%;              /* Gold/Yellow #FFC421 */
  --primary-foreground: 220 20% 10%;   /* Dark text on gold */

  --secondary: 220 25% 20%;            /* Dark slate for secondary elements */
  --secondary-foreground: 0 0% 95%;    /* Light text on secondary */

  --accent: 45 90% 55%;                /* Slightly lighter gold for accents */
  --accent-foreground: 220 20% 10%;    /* Dark text on accent */

  --muted: 220 15% 18%;                /* Muted dark */
  --muted-foreground: 220 10% 60%;     /* Muted text */

  --card: 220 18% 16%;                 /* Card background */
  --card-foreground: 0 0% 95%;         /* Card text */

  --popover: 220 18% 16%;              /* Popover background */
  --popover-foreground: 0 0% 95%;      /* Popover text */

  --destructive: 0 72% 51%;            /* Red for destructive */
  --destructive-foreground: 0 0% 100%; /* White on red */

  --input: 220 15% 25%;                /* Input background */
  --border: 220 15% 25%;               /* Border color */
  --ring: 45 100% 50%;                 /* Focus ring (gold) */

  /* Sidebar - slightly different shade for depth */
  --sidebar-background: 220 20% 11%;
  --sidebar-foreground: 0 0% 95%;
  --sidebar-primary: 45 100% 50%;
  --sidebar-primary-foreground: 220 20% 10%;
  --sidebar-accent: 220 20% 18%;
  --sidebar-accent-foreground: 0 0% 95%;
  --sidebar-border: 220 15% 20%;
  --sidebar-ring: 45 100% 50%;

  --radius: 0.5rem;
}
```

### Phase 2: Update Theme Metadata

**File**: `frontend/lib/data/themes.ts`

Add "hatchbridge" to the dark themes list so the system handles it correctly.

### Phase 3: Update Login Page Header

**File**: `frontend/pages/index.vue`

Update the header styling to use the dark background with gold accents:
- Dark background header
- White/gold text
- Gold accent buttons

### Phase 4: Update Default Layout Header

**File**: `frontend/layouts/default.vue`

Ensure the header bar uses the proper dark styling with gold accents.

### Phase 5: Update Logo Component (if needed)

**File**: `frontend/components/App/Logo.vue`

Ensure the SVG logo displays well on dark backgrounds (the current logo has yellow and black elements which should work).

### Phase 6: Update Button Styling

**File**: `frontend/components/ui/button/index.ts`

The default button variant should work with the new theme since it uses CSS variables. May need to ensure the "secondary" variant also looks good.

### Phase 7: Additional Component Updates

Review and update any hardcoded colors in:
- Cards and containers
- Form elements
- Navigation items
- Modal dialogs

---

## Specific File Changes

### 1. `frontend/assets/css/main.css`
- [ ] Update `:root, .homebox, .theme-hatchbridge` with new dark color scheme
- [ ] Ensure all color variables follow the Hatchbridge brand

### 2. `frontend/lib/data/themes.ts`
- [ ] Add "hatchbridge" to `darkThemes` array

### 3. `frontend/pages/index.vue`
- [ ] Update header color handling for hatchbridge theme
- [ ] Ensure gold/yellow buttons render correctly
- [ ] Update any hardcoded colors

### 4. `frontend/layouts/default.vue`
- [ ] Update header/sidebar styling for dark theme
- [ ] Ensure gold accents are prominent

### 5. `frontend/components/App/HeaderDecor.vue`
- [ ] Review and update decorative elements

### 6. `frontend/components/App/ThemePicker.vue`
- [ ] Ensure hatchbridge appears correctly in theme picker

---

## Color Reference

| Element | HSL | Hex (Approx) |
|---------|-----|--------------|
| Background | 220 15% 13% | #1e2128 |
| Background Accent | 220 18% 10% | #181b22 |
| Gold Primary | 45 100% 50% | #FFC421 |
| Foreground | 0 0% 98% | #FAFAFA |
| Muted | 220 15% 18% | #282d36 |
| Card | 220 18% 16% | #242831 |
| Border | 220 15% 25% | #363d4a |

---

## Testing Checklist

- [ ] Login page displays correctly with dark theme
- [ ] Dashboard/home page looks professional
- [ ] Sidebar navigation is clear and accessible
- [ ] Cards and forms are readable
- [ ] Buttons have proper contrast
- [ ] All text is readable
- [ ] Logo displays correctly
- [ ] Theme can be switched without issues
- [ ] Mobile responsive design works
