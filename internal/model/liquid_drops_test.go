package model

import (
	"testing"
	"time"
)

func TestPACCommit_ToLiquid(t *testing.T) {
	timestamp := time.Date(2026, 1, 7, 12, 0, 0, 0, time.UTC)
	commit := &PACCommit{
		SHA:       "abc123def456789",
		ShortSHA:  "abc123d",
		Header:    "Fix critical bug",
		Message:   "Fix critical bug\n\nDetailed explanation here.",
		Body:      "\nDetailed explanation here.",
		Timestamp: timestamp,
	}

	liquid := commit.ToLiquid()

	if liquid["sha"] != "abc123def456789" {
		t.Errorf("expected sha 'abc123def456789', got %v", liquid["sha"])
	}

	if liquid["shortsha"] != "abc123d" {
		t.Errorf("expected shortsha 'abc123d', got %v", liquid["shortsha"])
	}

	if liquid["header"] != "Fix critical bug" {
		t.Errorf("expected header 'Fix critical bug', got %v", liquid["header"])
	}

	if liquid["message"] != "Fix critical bug\n\nDetailed explanation here." {
		t.Errorf("unexpected message: %v", liquid["message"])
	}

	if liquid["time"] != timestamp {
		t.Errorf("expected time %v, got %v", timestamp, liquid["time"])
	}
}

func TestPACTask_ToLiquid(t *testing.T) {
	task := NewPACTask("JIRA-123")
	task.AddLabel("bug")
	task.AddLabel("critical")
	task.AddCommit(&PACCommit{SHA: "commit1", ShortSHA: "c1", Header: "First"})
	task.AddCommit(&PACCommit{SHA: "commit2", ShortSHA: "c2", Header: "Second"})
	task.Data = map[string]any{"status": "done"}
	task.Attributes["custom"] = "value"

	liquid := task.ToLiquid()

	if liquid["task_id"] != "JIRA-123" {
		t.Errorf("expected task_id 'JIRA-123', got %v", liquid["task_id"])
	}

	commits, ok := liquid["commits"].([]map[string]any)
	if !ok {
		t.Fatal("expected commits to be []map[string]any")
	}
	if len(commits) != 2 {
		t.Errorf("expected 2 commits, got %d", len(commits))
	}

	labels, ok := liquid["label"].([]string)
	if !ok {
		t.Fatal("expected label to be []string")
	}
	if len(labels) != 2 {
		t.Errorf("expected 2 labels, got %d", len(labels))
	}

	data, ok := liquid["data"].(map[string]any)
	if !ok {
		t.Fatal("expected data to be map[string]any")
	}
	if data["status"] != "done" {
		t.Errorf("expected data.status 'done', got %v", data["status"])
	}

	attrs, ok := liquid["attributes"].(map[string]any)
	if !ok {
		t.Fatal("expected attributes to be map[string]any")
	}
	if attrs["custom"] != "value" {
		t.Errorf("expected attributes.custom 'value', got %v", attrs["custom"])
	}
}

func TestPACTaskCollection_ToLiquid(t *testing.T) {
	collection := NewPACTaskCollection()

	// Add two referenced tasks
	task1 := collection.FindOrCreate("JIRA-123")
	task1.AddCommit(&PACCommit{SHA: "c1", Header: "Task 1 commit"})

	task2 := collection.FindOrCreate("JIRA-456")
	task2.AddCommit(&PACCommit{SHA: "c2", Header: "Task 2 commit"})

	// Add unreferenced commits
	unrefTask := collection.FindOrCreate("")
	unrefTask.AddCommit(&PACCommit{SHA: "u1", Header: "Unreferenced 1"})
	unrefTask.AddCommit(&PACCommit{SHA: "u2", Header: "Unreferenced 2"})

	liquid := collection.ToLiquid()

	referenced, ok := liquid["referenced"].([]map[string]any)
	if !ok {
		t.Fatal("expected referenced to be []map[string]any")
	}
	if len(referenced) != 2 {
		t.Errorf("expected 2 referenced tasks, got %d", len(referenced))
	}

	unreferenced, ok := liquid["unreferenced"].([]map[string]any)
	if !ok {
		t.Fatal("expected unreferenced to be []map[string]any")
	}
	if len(unreferenced) != 2 {
		t.Errorf("expected 2 unreferenced commits, got %d", len(unreferenced))
	}

	// Verify unreferenced commits have expected structure
	if unreferenced[0]["sha"] != "u1" {
		t.Errorf("expected first unreferenced sha 'u1', got %v", unreferenced[0]["sha"])
	}
}

func TestPACCommitCollection_ToLiquidCommits(t *testing.T) {
	collection := NewPACCommitCollection()
	collection.Add(&PACCommit{SHA: "c1", Header: "First"})
	collection.Add(&PACCommit{SHA: "c2", Header: "Second"})

	liquid := collection.ToLiquidCommits()

	if len(liquid) != 2 {
		t.Errorf("expected 2 commits, got %d", len(liquid))
	}

	if liquid[0]["sha"] != "c1" {
		t.Errorf("expected first sha 'c1', got %v", liquid[0]["sha"])
	}

	if liquid[1]["header"] != "Second" {
		t.Errorf("expected second header 'Second', got %v", liquid[1]["header"])
	}
}
