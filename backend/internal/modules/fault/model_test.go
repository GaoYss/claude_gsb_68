package fault

import "testing"

// TestCanTransitToRejectsSelfLoop 任何状态都不允许原地不动,
// 同状态重复提交必须被状态机拒绝, 避免重复刷新处置时间。
func TestCanTransitToRejectsSelfLoop(t *testing.T) {
	for _, status := range Statuses() {
		if canTransitTo(status, status) {
			t.Fatalf("状态 %s 不允许迁移到自身", status)
		}
	}
}

func TestCanTransitToLegalEdges(t *testing.T) {
	legal := [][2]string{
		{StatusPending, StatusProcessing},
		{StatusPending, StatusClosed},
		{StatusProcessing, StatusRepaired},
		{StatusProcessing, StatusClosed},
		{StatusRepaired, StatusClosed},
		{StatusRepaired, StatusProcessing}, // 返修
	}
	for _, edge := range legal {
		if !canTransitTo(edge[0], edge[1]) {
			t.Fatalf("合法迁移 %s -> %s 被误拒", edge[0], edge[1])
		}
	}
}

func TestCanTransitToIllegalEdges(t *testing.T) {
	illegal := [][2]string{
		{StatusPending, StatusRepaired}, // 不能跳过维修直接修复
		{StatusRepaired, StatusPending},
		{StatusClosed, StatusPending},
		{StatusClosed, StatusProcessing},
		{StatusClosed, StatusRepaired},
		{StatusClosed, StatusClosed},
	}
	for _, edge := range illegal {
		if canTransitTo(edge[0], edge[1]) {
			t.Fatalf("非法迁移 %s -> %s 被误放", edge[0], edge[1])
		}
	}
}

// TestCanStartRepair 维修中允许再次开工(状态不变的新事件), 已修复/已关闭则不允许。
func TestCanStartRepair(t *testing.T) {
	if !canStartRepair(StatusPending) {
		t.Fatal("待处理故障应允许首次开工")
	}
	if !canStartRepair(StatusProcessing) {
		t.Fatal("维修中应允许再次开工(二次维修)")
	}
	if canStartRepair(StatusRepaired) {
		t.Fatal("已修复故障不应通过开工校验, 返修需先走 repaired -> processing 迁移")
	}
	if canStartRepair(StatusClosed) {
		t.Fatal("已关闭故障不允许开工")
	}
}
