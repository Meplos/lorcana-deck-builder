// Package Ink: color.go utils function to convert and interact with ink color
package ink

type (
	InkColor string
	InkMask  int64
)

const (
	AMBER    InkMask = 1
	AMETHYST InkMask = 2
	EMERALD  InkMask = 4
	RUBY     InkMask = 8
	SAPPHIRE InkMask = 16
	STEEL    InkMask = 32
)

func GetInkStrings(mask InkMask) []InkColor {
	colors := make([]InkColor, 0)
	if mask&AMBER != 0 {
		colors = append(colors, "amber")
	}
	if mask&AMETHYST != 0 {
		colors = append(colors, "amethyst")
	}
	if mask&EMERALD != 0 {
		colors = append(colors, "emerald")
	}
	if mask&RUBY != 0 {
		colors = append(colors, "ruby")
	}
	if mask&SAPPHIRE != 0 {
		colors = append(colors, "sapphire")
	}
	if mask&STEEL != 0 {
		colors = append(colors, "steel")
	}

	return colors
}
