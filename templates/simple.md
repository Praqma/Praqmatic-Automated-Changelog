# {{.properties.title}} for {{.properties.product}}
## All Referenced Tasks
{{range referencedTasks .tasks}}
### {{.TaskID}}
{{range .Commits}}
- {{.ShortSHA}}: {{.Header}}
{{end}}{{end}}
## Unspecified Commits
{{range unreferencedCommits .tasks}}
- {{.ShortSHA}}: {{.Header}}
{{end}}