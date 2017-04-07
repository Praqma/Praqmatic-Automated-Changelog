# Migrating from 2.x to 3.x

With the introduction of PAC 3.0, a couple of things have changed:

1. The previously deprecated `--sha` parameter has been completely removed, as `from` provides the exact same functionality.
2. All `date`-related functionality has been removed, as it was imprecise and difficult to use.

This means that the only remaining ways of querying is `from` and `--to` as well as `from-latest-tag`, greatly simplifying use of pac.
With the removal of `date`, there is no longer any reason to have `date_template` defined in settings files. When upgrading to 3.0, it is recommended to remove this parameter from settings files, or at least make sure not to add it to any new files. It is not a breaking change, as the parameter is just unused in the future.

## Old format

```YAML
:general:
  date_template: "%Y-%m-%d"
  changelog_name: "changelog"
  changelog_formats:
    - "html"
    - "pdf"
  changelog_css:
  verbose: false
[...]
```

## New format (date_template removed)

```YAML
:general:
  changelog_name: "changelog"
  changelog_formats:
    - "html"
    - "pdf"
  changelog_css:
  verbose: false
[...]
```