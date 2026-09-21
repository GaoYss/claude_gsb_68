package fault

import "testing"

// TestCanTransitTo 验证状态机流转表, 同状态原地不动不属于合法流转。
func TestCanTransitTo(t *testing.T) {
	legal := map[string][]string{
		StatusPending:    {StatusProcessing, StatusClosed},
		StatusProcessing: {StatusRepaired, StatusClosed},
		StatusRepaired:   {StatusClosed, StatusProcessing},
		StatusClosed:     {},
	}

	for _, from := range Statuses() {
		for _, to := range Statuses() {
			got := canTransitTo(from, to)
			want := contains(legal[from], to)
			if got != want {
				t.Errorf("canTransitTo(%s, %s) = %v, 期望 %v", from, to, got, want)
			}
			if from == to && got {
				t.Errorf("canTransitTo(%s, %s) 不应允许原地不动", from, to)
			}
		}
	}
}

func contains(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}
