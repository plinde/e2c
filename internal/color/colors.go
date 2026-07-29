// SPDX-FileCopyrightText: Copyright (C) Nicolas Lamirault <nicolas.lamirault@gmail.com>
// SPDX-License-Identifier: Apache-2.0

package color

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// DefaultTheme is the theme used when none is configured.
const DefaultTheme = "nord"

// Colors represents the colors used in the application
type Colors struct {
	// Basic colors
	Background tcell.Color
	Foreground tcell.Color
	Border     tcell.Color
	Title      tcell.Color
	Selected   tcell.Color

	// UI element colors
	HeaderFg tcell.Color
	HeaderBg tcell.Color

	// Status colors
	Running tcell.Color
	Stopped tcell.Color
	Pending tcell.Color
	Error   tcell.Color

	// Other colors
	Highlight tcell.Color
	Secondary tcell.Color
}

// Themes holds the built-in color schemes, keyed by the name accepted by
// --theme and by the ui.theme config key.
var Themes = map[string]Colors{
	// https://www.nordtheme.com/
	"nord": {
		Background: tcell.GetColor("#2E3440"), // Primary background
		Foreground: tcell.GetColor("#D8DEE9"), // Primary foreground
		Border:     tcell.GetColor("#81A1C1"), // Normal blue
		Title:      tcell.GetColor("#88C0D0"), // Normal cyan
		Selected:   tcell.GetColor("#3B4252"), // Normal black
		HeaderFg:   tcell.GetColor("#ECEFF4"), // Bright white
		HeaderBg:   tcell.GetColor("#4C566A"), // Bright black
		Running:    tcell.GetColor("#A3BE8C"), // Normal green
		Stopped:    tcell.GetColor("#BF616A"), // Normal red
		Pending:    tcell.GetColor("#EBCB8B"), // Normal yellow
		Error:      tcell.GetColor("#BF616A"), // Normal red
		Highlight:  tcell.GetColor("#EBCB8B"), // Normal yellow
		Secondary:  tcell.GetColor("#81A1C1"), // Normal blue
	},

	// https://github.com/morhetz/gruvbox — dark, using the bright accent
	// variants for legibility against bg0, mirroring the Nord mapping above.
	"gruvbox": {
		Background: tcell.GetColor("#282828"), // bg0
		Foreground: tcell.GetColor("#ebdbb2"), // fg1
		Border:     tcell.GetColor("#83a598"), // Bright blue
		Title:      tcell.GetColor("#8ec07c"), // Bright aqua
		Selected:   tcell.GetColor("#504945"), // bg2
		HeaderFg:   tcell.GetColor("#fbf1c7"), // fg0
		HeaderBg:   tcell.GetColor("#665c54"), // bg3
		Running:    tcell.GetColor("#b8bb26"), // Bright green
		Stopped:    tcell.GetColor("#fb4934"), // Bright red
		Pending:    tcell.GetColor("#fabd2f"), // Bright yellow
		Error:      tcell.GetColor("#fb4934"), // Bright red
		Highlight:  tcell.GetColor("#fabd2f"), // Bright yellow
		Secondary:  tcell.GetColor("#83a598"), // Bright blue
	},
}

// AppColors is the active color scheme.
var AppColors = Themes[DefaultTheme]

// ThemeNames returns the built-in theme names, sorted.
func ThemeNames() []string {
	names := make([]string, 0, len(Themes))
	for name := range Themes {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// SetTheme replaces the active color scheme with the named built-in theme.
// An empty name selects DefaultTheme. Unknown names are an error rather than a
// silent fallback, so a typo in the config surfaces immediately.
func SetTheme(name string) error {
	if name == "" {
		name = DefaultTheme
	}
	theme, ok := Themes[strings.ToLower(name)]
	if !ok {
		return fmt.Errorf("unknown theme %q (available: %s)", name, strings.Join(ThemeNames(), ", "))
	}
	AppColors = theme
	return nil
}

// fields maps a config key to the corresponding field of AppColors.
func fields() map[string]*tcell.Color {
	return map[string]*tcell.Color{
		"background": &AppColors.Background,
		"foreground": &AppColors.Foreground,
		"border":     &AppColors.Border,
		"title":      &AppColors.Title,
		"selected":   &AppColors.Selected,
		"header_fg":  &AppColors.HeaderFg,
		"header_bg":  &AppColors.HeaderBg,
		"running":    &AppColors.Running,
		"stopped":    &AppColors.Stopped,
		"pending":    &AppColors.Pending,
		"error":      &AppColors.Error,
		"highlight":  &AppColors.Highlight,
		"secondary":  &AppColors.Secondary,
	}
}

// ColorNames returns the overridable color keys, sorted.
func ColorNames() []string {
	f := fields()
	names := make([]string, 0, len(f))
	for name := range f {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// ApplyOverrides overrides individual colors of the active theme. Values are
// anything tcell understands — "#rrggbb" or a color name such as "red".
// Unknown keys and unparseable values are errors, again so config typos are not
// silently ignored.
func ApplyOverrides(overrides map[string]string) error {
	f := fields()
	// Sorted for a deterministic error on multiple bad keys.
	keys := make([]string, 0, len(overrides))
	for key := range overrides {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		target, ok := f[strings.ToLower(key)]
		if !ok {
			return fmt.Errorf("unknown color %q (available: %s)", key, strings.Join(ColorNames(), ", "))
		}
		value := overrides[key]
		parsed := tcell.GetColor(value)
		if parsed == tcell.ColorDefault {
			return fmt.Errorf("invalid color value %q for %q (expected #rrggbb or a color name)", value, key)
		}
		*target = parsed
	}
	return nil
}

// InitializeColors applies the application colors to tview components
func InitializeColors() {
	// Apply colors to tview global styles
	tview.Styles.PrimitiveBackgroundColor = AppColors.Background
	tview.Styles.ContrastBackgroundColor = AppColors.HeaderBg
	tview.Styles.MoreContrastBackgroundColor = AppColors.Selected
	tview.Styles.BorderColor = AppColors.Border
	tview.Styles.TitleColor = AppColors.Title
	tview.Styles.GraphicsColor = AppColors.Border
	tview.Styles.PrimaryTextColor = AppColors.Foreground
	tview.Styles.SecondaryTextColor = AppColors.Secondary
	tview.Styles.TertiaryTextColor = AppColors.Highlight
	tview.Styles.InverseTextColor = AppColors.HeaderFg
	tview.Styles.ContrastSecondaryTextColor = AppColors.Highlight
}
