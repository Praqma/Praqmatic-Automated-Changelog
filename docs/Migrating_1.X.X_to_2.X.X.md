## Migrating from 1.x to 2.x

With the introduction of PAC 2.0, the settings file formatting has changed.
The change isn't backwards compatible, so you will need to edit your settings file when upgrading from 1.x.

In the following example, we upgrade a simple 1.x settings file:

#### Old format

```YAML
:general:
  date_template: "%Y-%m-%d"
  changelog_name: "changelog"
  changelog_formats:
    - "html"
    - "pdf"
  changelog_css:
  verbose: false

:none:
  regex:
    - '/(JENKINS-[0-9]+)/i'

:vcs:
  type: git
  usr:
  pwd:
  repo_location: "."
```
#### New format

```YAML
:general:
  date_template: '%Y-%m-%d'
  :strict: false

:templates:
  - { location: templates/default.md, output: changelog.md }
  - { location: templates/default_html.html, pdf: true, output: changelog.html }

:task_systems:
  -
    :name: none
    :regex:
      - { pattern: '/(JENKINS-[0-9]+)', label: none }
    :delimiter: '/,|\s/'

:vcs:
  :type: git
  :repo_location: '.'
```

## Migration steps

The following points cover all the settings file sections and the changes required to upgrade them.

### :general:

  1. Remove `changelog_name` and `changelog_formats`. These setting are now part of the `:templates:` section.
  2. Remove `changelog_css:` and `verbose:` as they have been deprecated.

### :templates:

Previous attributes `changelog_name`, `changelog_css` and `changelog_formats` are all deprecated and replaced with templates.
You need to have at least one template specified to produce an output.

To generate three output files in 1.x (`changelog.md`, `changelog.html` and `changelog.pdf`), you would have specified the following options in your settings file: 

```YAML
:general:
	changelog_name: "Changelog"
    changelog_formats:
      - "html"
      - "pdf"
```

To produce the same output files in the new format, replace the old configuration with templates. Add the following as your `:templates:` section

```YAML
:templates:
  - { location: templates/default.md, output: Changelog.md }
  - { location: templates/default_html.html, pdf: true, output: Changelog.html }
```

Don't forget to create the Liquid templates defined in the `location:` attribute on the template items. For examples of Liquid templates, take a look at the examples included in the `templates` folder of this project. More information regarding templates can also be found in the project's `README` file.

### :task_systems: 

If you used the task systems `none` or `trac`, move them under the `:task_system:` section using the list notation demonstrated below:

```YAML
:task_systems:
	-
		:name: '[trac/none]'
```

Next, add a `:regex:` section under your task system. Move your regex into the new section as a `pattern` and asssign a grouping `label` to it.

```YAML
		:regex:
		  - { pattern: '/[your pattern]', label: [matching_issues] }
```

_Note:_ Previously a capture group named `<id>` was required in your regex. This is no longer the case. Now the first capture group is used and one regex can have multiple matches on the same commit.

### :vcs: 

Both `usr` and `pw` have been removed.
