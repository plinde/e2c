// SPDX-FileCopyrightText: Copyright (C) Nicolas Lamirault <nicolas.lamirault@gmail.com>
// SPDX-License-Identifier: Apache-2.0

package color

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

// restore resets the package-level AppColors after a test mutates it.
func restore(t *testing.T) {
	t.Helper()
	saved := AppColors
	t.Cleanup(func() { AppColors = saved })
}

func TestDefaultThemeExists(t *testing.T) {
	if _, ok := Themes[DefaultTheme]; !ok {
		t.Fatalf("DefaultTheme %q is not present in Themes", DefaultTheme)
	}
}

func TestSetTheme(t *testing.T) {
	tests := []struct {
		name    string
		theme   string
		wantBg  tcell.Color
		wantErr bool
	}{
		{name: "empty selects default", theme: "", wantBg: Themes[DefaultTheme].Background},
		{name: "nord", theme: "nord", wantBg: tcell.GetColor("#2E3440")},
		{name: "gruvbox", theme: "gruvbox", wantBg: tcell.GetColor("#282828")},
		{name: "case insensitive", theme: "GruvBox", wantBg: tcell.GetColor("#282828")},
		{name: "unknown is an error", theme: "solarized", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			restore(t)
			err := SetTheme(tt.theme)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("SetTheme(%q) = nil, want error", tt.theme)
				}
				return
			}
			if err != nil {
				t.Fatalf("SetTheme(%q) = %v, want nil", tt.theme, err)
			}
			if AppColors.Background != tt.wantBg {
				t.Errorf("Background = %v, want %v", AppColors.Background, tt.wantBg)
			}
		})
	}
}

func TestApplyOverrides(t *testing.T) {
	t.Run("overrides a single color", func(t *testing.T) {
		restore(t)
		if err := SetTheme("nord"); err != nil {
			t.Fatalf("SetTheme: %v", err)
		}
		if err := ApplyOverrides(map[string]string{"running": "#b8bb26"}); err != nil {
			t.Fatalf("ApplyOverrides: %v", err)
		}
		if got, want := AppColors.Running, tcell.GetColor("#b8bb26"); got != want {
			t.Errorf("Running = %v, want %v", got, want)
		}
		// Untouched fields keep the theme value.
		if got, want := AppColors.Background, tcell.GetColor("#2E3440"); got != want {
			t.Errorf("Background = %v, want %v", got, want)
		}
	})

	t.Run("nil map is a no-op", func(t *testing.T) {
		restore(t)
		if err := ApplyOverrides(nil); err != nil {
			t.Fatalf("ApplyOverrides(nil) = %v, want nil", err)
		}
	})

	t.Run("unknown key is an error", func(t *testing.T) {
		restore(t)
		if err := ApplyOverrides(map[string]string{"bakground": "#000000"}); err == nil {
			t.Fatal("ApplyOverrides with unknown key = nil, want error")
		}
	})

	t.Run("invalid value is an error", func(t *testing.T) {
		restore(t)
		if err := ApplyOverrides(map[string]string{"running": "not-a-color"}); err == nil {
			t.Fatal("ApplyOverrides with invalid value = nil, want error")
		}
	})
}

func TestColorNamesCoverEveryField(t *testing.T) {
	// fields() is the override surface; keep it in sync with the Colors struct.
	if got, want := len(ColorNames()), 13; got != want {
		t.Errorf("ColorNames() has %d entries, want %d — did Colors gain a field?", got, want)
	}
}
