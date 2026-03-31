# Design System Document: High-End Editorial Insurance Experience

## 1. Overview & Creative North Star: "The Digital Concierge"
This design system moves away from the "utility-first" aesthetic common in the insurance industry toward a **High-End Editorial** experience. Our Creative North Star is **"The Digital Concierge."**

In high-stress scenarios like auto claims, the UI must act as a calm, authoritative guide. We achieve this by breaking the rigid, boxed-in "template" look. Instead, we utilize intentional asymmetry, significant negative space, and a sophisticated layering of surfaces. By rejecting traditional 1px borders and embracing tonal depth, we create an environment that feels bespoke, premium, and inherently trustworthy.

---

## 2. Colors & Surface Philosophy
The palette is rooted in a deep navy and slate foundation, punctuated by a "Success Green" that acts as a beacon of resolution.

### The "No-Line" Rule
To maintain a high-end feel, **1px solid borders for sectioning are strictly prohibited.** Boundaries must be defined through:
- **Background Color Shifts:** Use `surface-container-low` (#f3f4f0) against the main `surface` (#f9faf5).
- **Tonal Transitions:** Defining functional areas by nesting different surface tiers rather than drawing lines around them.

### Surface Hierarchy & Nesting
Treat the UI as a series of physical layers—like stacked sheets of fine stationery.
- **Base Layer:** `surface` (#f9faf5) – The canvas for the entire application.
- **Secondary Areas:** `surface-container` (#edeeea) – Used for sidebars or secondary navigation.
- **Interactive Elements:** `surface-container-lowest` (#ffffff) – Reserved for the most important cards and inputs to make them "pop" against the off-white background.

### Signature Textures & Glassmorphism
- **The Glass Rule:** For floating modals or navigation bars, use `surface` at 80% opacity with a `backdrop-filter: blur(20px)`. This integrates the component into the environment.
- **The Narrative Gradient:** For primary CTAs and Hero sections, use a subtle linear gradient transitioning from `primary` (#000000) to `primary_container` (#101b30) at a 135-degree angle. This adds "soul" and prevents the deep navy from feeling flat or "dead."

---

## 3. Typography: Authoritative Clarity
We pair **Manrope** (Display/Headline) with **Inter** (Body/Labels) to balance editorial sophistication with high-stress readability.

- **Display (Manrope):** Large, bold scales (3.5rem - 2.25rem) are used for "Atmospheric Statements"—welcoming the user or stating a claim status. Use `on_surface` (#1a1c1a) with -0.02em letter spacing for a premium "ink-on-paper" feel.
- **Headlines (Manrope):** Used for section headers. These provide the "Voice of Authority."
- **Body (Inter):** Set in `on_surface_variant` (#44474c) for secondary info and `on_surface` for primary text. Inter’s tall x-height ensures maximum legibility during claim filing.
- **Labels (Inter):** Small, all-caps or high-weight labels using `secondary` (#47607e) to categorize information without cluttering the visual field.

---

## 4. Elevation & Depth
In this system, depth is a functional tool, not a decoration.

- **Tonal Layering:** Place a `surface-container-lowest` (#ffffff) card on a `surface-container-low` (#f3f4f0) background. This creates a natural, soft lift that signals "interactivity" without the visual noise of a shadow.
- **Ambient Shadows:** When a floating effect is required (e.g., a "Submit Claim" FAB), use an extra-diffused shadow: `box-shadow: 0 12px 32px rgba(16, 27, 48, 0.06)`. Note the use of a navy-tinted shadow rather than pure black.
- **The Ghost Border:** If accessibility requires a stroke (e.g., high-contrast mode), use `outline_variant` (#c4c6cc) at **15% opacity**. It should be felt, not seen.

---

## 5. Components

### Buttons: The Elegant Call-to-Action
- **Primary:** Gradient fill (`primary` to `primary_container`). Border radius `md` (0.375rem). Text is `on_primary` (#ffffff).
- **Secondary:** `surface-container-highest` (#e2e3df) background with `on_surface` text. No border.
- **Tertiary:** Text-only using `secondary` (#47607e) with a subtle underline on hover.

### Inputs: Crisp & Clear
- **Field:** Background `surface-container-lowest` (#ffffff). Border: Ghost Border (15% `outline`).
- **Focus State:** 2px solid `primary` (#000000). The sharp transition from a ghost border to a solid primary line provides a "high-tech" responsive feel.

### Cards & Lists: The "No-Divider" Rule
- **Cards:** Forbid the use of divider lines within cards. Use `spacing-6` (1.5rem) to separate header from body text.
- **List Items:** Use a subtle background hover state (`surface-container-high`) rather than a line between items. This keeps the interface feeling "open" and breathable.

### Signature Component: The "Status Ribbon"
For insurance claims, use a `tertiary_fixed` (#b1f0ce) pill with `on_tertiary_fixed_variant` (#0e5138) text. Position it asymmetrically in the top-right of a card to break the grid and draw the eye immediately to the claim's progress.

---

## 6. Do’s and Don’ts

### Do
- **Do** use `spacing-16` (4rem) and `spacing-20` (5rem) for section margins to create an editorial, "breathable" layout.
- **Do** use asymmetric image placements (e.g., an auto-photo offset from its container) to feel more like a premium magazine than a database.
- **Do** ensure all "Success" states use the `tertiary` (#000000) and `tertiary_fixed` tokens to provide a calming, reliable psychological cue.

### Don't
- **Don't** use 100% opaque `outline` (#74777d) for borders. It creates "visual cages" that increase user anxiety.
- **Don't** use standard "Drop Shadows." Only use the ambient, tinted navy shadows specified in Section 4.
- **Don't** crowd the interface. If the screen feels full, increase the `spacing` tokens. High-end design requires the luxury of space.