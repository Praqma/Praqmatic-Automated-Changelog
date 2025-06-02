# {{.properties.title}} for {{.properties.product}}
## Tasks by Author
{{range $author, $tasks := groupByAuthor (referencedTasks .tasks)}}
### {{$author}}
{{range $tasks}}
#### {{.TaskID}} - {{.Title}}
- State: {{.State}}
- Labels: {{joinLabels .Labels ", "}}
Commits:
{{range .Commits}}
- {{commitSHA . 7}}: {{commitHeader .}}
{{end}}{{end}}{{end}}
## All Referenced Tasks
{{range referencedTasks .tasks}}
### {{.TaskID}} - {{.Title}}
{{range .Commits}}
- {{.ShortSHA}}: {{.Header}}
{{end}}{{end}}
## Unspecified Commits
{{range unreferencedCommits .tasks}}
- {{.ShortSHA}}: {{.Header}}
{{end}}
## Statistics
{{$stats := taskStats (allTasks .tasks)}}
- Total Tasks: {{$stats.total}}
- Open: {{$stats.open}} | Closed: {{$stats.closed}}
- Issues: {{$stats.issues}} | Pull Requests: {{$stats.pull_requests}}
## Closed Tasks
{{range filterByState (referencedTasks .tasks) "closed"}}
- [{{.TaskID}}]({{issueURL .}}) - {{.Title}}
{{end}}
