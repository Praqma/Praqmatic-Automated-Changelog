package model

import (
	"testing"
)

func TestPACCommitCollection_Add(t *testing.T) {
	collection := NewPACCommitCollection()

	commit := &PACCommit{
		SHA:      "abc123def456",
		ShortSHA: "abc123d",
		Header:   "Fix bug",
		Message:  "Fix bug\n\nDetailed description",
	}

	collection.Add(commit)

	if collection.Count() != 1 {
		t.Errorf("expected count 1, got %d", collection.Count())
	}

	if collection.Commits[0] != commit {
		t.Error("commit not added correctly")
	}
}

func TestPACCommitCollection_CountWithAndWithout(t *testing.T) {
	collection := NewPACCommitCollection()

	// Add 3 referenced commits and 2 unreferenced
	for i := 0; i < 3; i++ {
		commit := &PACCommit{SHA: "ref" + string(rune('0'+i)), Referenced: true}
		collection.Add(commit)
	}
	for i := 0; i < 2; i++ {
		commit := &PACCommit{SHA: "unref" + string(rune('0'+i)), Referenced: false}
		collection.Add(commit)
	}

	if collection.Count() != 5 {
		t.Errorf("expected total count 5, got %d", collection.Count())
	}

	if collection.CountWith() != 3 {
		t.Errorf("expected CountWith 3, got %d", collection.CountWith())
	}

	if collection.CountWithout() != 2 {
		t.Errorf("expected CountWithout 2, got %d", collection.CountWithout())
	}
}

func TestPACCommitCollection_Health(t *testing.T) {
	tests := []struct {
		name           string
		referenced     int
		unreferenced   int
		expectedHealth float64
	}{
		{
			name:           "empty collection",
			referenced:     0,
			unreferenced:   0,
			expectedHealth: 100.0,
		},
		{
			name:           "all referenced",
			referenced:     10,
			unreferenced:   0,
			expectedHealth: 100.0,
		},
		{
			name:           "none referenced",
			referenced:     0,
			unreferenced:   10,
			expectedHealth: 0.0,
		},
		{
			name:           "half referenced",
			referenced:     5,
			unreferenced:   5,
			expectedHealth: 50.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collection := NewPACCommitCollection()
			for i := 0; i < tt.referenced; i++ {
				collection.Add(&PACCommit{Referenced: true})
			}
			for i := 0; i < tt.unreferenced; i++ {
				collection.Add(&PACCommit{Referenced: false})
			}

			if health := collection.Health(); health != tt.expectedHealth {
				t.Errorf("expected health %.1f, got %.1f", tt.expectedHealth, health)
			}
		})
	}
}

func TestPACTask_Labels(t *testing.T) {
	task := NewPACTask("JIRA-123")

	task.AddLabel("bug")
	task.AddLabel("feature")

	if !task.HasLabel("bug") {
		t.Error("expected task to have 'bug' label")
	}

	if !task.HasLabel("feature") {
		t.Error("expected task to have 'feature' label")
	}

	if task.HasLabel("enhancement") {
		t.Error("expected task NOT to have 'enhancement' label")
	}

	labels := task.GetLabels()
	if len(labels) != 2 {
		t.Errorf("expected 2 labels, got %d", len(labels))
	}

	// Test ClearLabels
	task.ClearLabels()
	if len(task.Labels) != 0 {
		t.Error("expected labels to be cleared")
	}
}

func TestPACTaskCollection_FindOrCreate(t *testing.T) {
	collection := NewPACTaskCollection()

	// Create new task
	task1 := collection.FindOrCreate("JIRA-123")
	if task1.TaskID != "JIRA-123" {
		t.Errorf("expected task ID 'JIRA-123', got %q", task1.TaskID)
	}

	// Find existing task
	task2 := collection.FindOrCreate("JIRA-123")
	if task1 != task2 {
		t.Error("expected FindOrCreate to return same task instance")
	}

	// Create another task
	task3 := collection.FindOrCreate("JIRA-456")
	if task3.TaskID != "JIRA-456" {
		t.Errorf("expected task ID 'JIRA-456', got %q", task3.TaskID)
	}

	if collection.Count() != 2 {
		t.Errorf("expected 2 tasks, got %d", collection.Count())
	}
}

func TestPACTaskCollection_ReferencedAndUnreferenced(t *testing.T) {
	collection := NewPACTaskCollection()

	// Add referenced tasks
	task1 := collection.FindOrCreate("JIRA-123")
	task1.AddCommit(&PACCommit{SHA: "commit1"})

	task2 := collection.FindOrCreate("JIRA-456")
	task2.AddCommit(&PACCommit{SHA: "commit2"})

	// Add unreferenced commits (empty task ID)
	unrefTask := collection.FindOrCreate("")
	unrefTask.AddCommit(&PACCommit{SHA: "commit3"})
	unrefTask.AddCommit(&PACCommit{SHA: "commit4"})

	// Test Referenced
	referenced := collection.Referenced()
	if len(referenced) != 2 {
		t.Errorf("expected 2 referenced tasks, got %d", len(referenced))
	}

	// Test Unreferenced
	unreferenced := collection.Unreferenced()
	if unreferenced == nil {
		t.Fatal("expected unreferenced task to exist")
	}
	if len(unreferenced.Commits.Commits) != 2 {
		t.Errorf("expected 2 unreferenced commits, got %d", len(unreferenced.Commits.Commits))
	}

	// Test UnreferencedCommits
	unrefCommits := collection.UnreferencedCommits()
	if len(unrefCommits) != 2 {
		t.Errorf("expected 2 unreferenced commits, got %d", len(unrefCommits))
	}
}

func TestPACTaskCollection_NoUnreferenced(t *testing.T) {
	collection := NewPACTaskCollection()

	// Only add referenced tasks
	collection.FindOrCreate("JIRA-123")

	unreferenced := collection.Unreferenced()
	if unreferenced != nil {
		t.Error("expected no unreferenced task")
	}

	unrefCommits := collection.UnreferencedCommits()
	if len(unrefCommits) != 0 {
		t.Error("expected empty unreferenced commits slice")
	}
}
