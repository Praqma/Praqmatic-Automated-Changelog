# PAC Changelog
{{range .tasks.referenced}}
## {{.TaskID}}
{{range .Commits}}
- {{.ShortSHA}}: {{.Header}}
{{end}}
{{end}}
## Unspecified
{{range .tasks.unreferenced}}
- {{.ShortSHA}}: {{.Header}}
{{end}}
