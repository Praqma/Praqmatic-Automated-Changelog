# Roadmap

All changes are to comform and support the following visions for PAC described in the section below.

The [original blue print](/docs/BluePrint.md) behind a pragmatic automated changelog also still serves as guidelines, especially statememts like these:

> An important aspect of a changelog is that it should be targeted to a specific groups of readers. It could be project managers or it could be end-user and customers. This might requires two very different worded documents, though they are in the end based on the exact same information from the changelog.

_so we are able to collect information from several systems, one for developers and maybe one with end user docs_.

> In a continuous delivery world we are always ready to release, and as a release also include a changelog, or a release note, and other kind of documentation, these documents must also always be available continuously. .. we recognize that the good release note requires some manual work ...

_so doing it continuously for every commit make missing work visible._



## Future epics

There is already ongoing work, or plans for the features and usage within the following areas. You will find other work accepted on our [github issues backlog](https://github.com/Praqma/Praqmatic-Automated-Changelog/issues?q=is%3Aopen+is%3Aissue+milestone%3Abacklog) as well, but these are our major areas of focus.


* **Multiple query support**: Changes are being discussed and planned, to allow to describe multiple queries when collecting information. Queries will be able to use the information collected from SCM, or from other queries.
* **More flexibility through plugin**: We are focusing on working toward more flexibility by supporting some kind of plugin architecture. This means ability to collect information by calling external scripts or support even more systems than the usual tasks systems.
* **Continues support for both running as a script or in a container**: Focus is still to support using PAC as a script going forward, so it can run a script from the repository. Support for distributing it as container or other packages also continues.
* **Cross-platform - more windows friendly-ness**: All changes have to be platform agnostic or work cross-platform. We should not only abstract by using containers, but must also support that PAC continues to work as a script in the runtime environment meaning any 3rd party dependencies like libraries used must be available (easily!) on major platforms (Windows, MAC, Linux)

_A general rule of thumb is we don't accept changes that break the above plans, or make the impossible or very difficult to implement._

## Architecture

The original architecture still gives the overal view and background for design decisions: [/docs/designview.png](/docs/designview.png).

**Not all details in the implementation honors this picture - thus our future plans are to conform.**

Key elements are:

* Module for parsing and working with configuration, currently integrated in the core
* VCS modules, using common interface
* Writer module, that writes the report, interface to the core
* Task modules, those that collects information based on SCM information (or other supplied look-up information).

_Interfaces_ are mostly common data structures used to communicate between the modules.


## Documentation

**Documentation is to be simplified by moving to use-case oriented documentation and illustration**: Today there is a large amount of text, which becomes more difficult to maintain properly so we move towards a simplied documentation setup that originates in use-cases and how to use PAC for them. Further many of the detailed writing are to be transformed or improved using drawing or pictures.

## Testing and quality

Today we have a fair amount of tests on different levels in PAC, from unit tests - that are not really unit tests, to functional tests and integration tests that interacts with real systems.

When changing code, we have the following goals:

* all new code must be tested in an automated reproducible manner
* changed existing code must have improved tests if coverage (on code level or use-case level) isn't good enough

The roadmap for testings is moving in direction of:

* **Test-cases on use-case level**: We want to track which use-cases are tested and what tests belong to what use-cases. Note the documentation roadmap that moves documentation towards being use-case oriented.
* **Focus on several layers of testing**: We will move towards the test strategy described below, when we work with tests.

### Test strategy

We will be doing fully automated and reproducible tests in the project, prioritized in the following order - at least one needs to apply:

* prefer unit testings if possible
* functional testings when possible, using moc data for queries are okay
* some real integration tests that needs to interact with live systems for stability
* build up a regression test suite, possible re-use of above tests, to ensure primary use-cases are always working

See also [developer information about tests](/docs/developer_info.md#tests)


## Use the roadmap

The roadmap should be live and changeable, but should set some future goals on a very high level. Though the roadmap have some visions, they might be up for discussion.

The roadmap should in the daily work be seen as:

* a guideline and check list for what changes we accept and what work can be put on our backlog
* information to any developer, that implementations and design decision should honor our future vision - or at least not make it harder to achive


_And example on the last one, is that we can not ever include any more non-cross-platform libraries, because we epic about being cross-platform in the future_.
