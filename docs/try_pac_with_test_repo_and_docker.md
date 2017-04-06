# Try PAC with a test repo and PAC docker image

Take PAC for a spin...

Follow these steps to try PAC with a test repository and the PAC docker image we supply to see a more elaborate example on how PAC works.

Unzip one of the git repositories we use for testing - it can be used to generate a simple changelog also.

* clone the [PAC Github repository](https://github.com/Praqma/Praqmatic-Automated-Changelog) to your local computer
* Unzip the file `test/resources/idReportTestRepository.zip` file from this repository to a folder on your computer
* Create a file `default_settings.yml` and paste the contents from below into this file. Put the file in the root of the extracted Git repository 

```
:general:
  :strict: true

:templates:
  - { location: /usr/src/app/templates/default_id_report.md, output: ids.md }
  - { location: /usr/src/app/templates/default.md, output: default.md }
  - { location: /usr/src/app/templates/default_html.html, pdf: true, output: default.html }

:task_systems:
  - 
    :name: none
    :regex:
      - { pattern: '/.*Issue:\s*(?<id>[\d+|[,|\s]]+).*?\n/im', label: none }
      - { pattern: '/.*Issue:\s*?(none).*?\n/im', label: none}
    :delimiter: '/,|\s/'
  
:vcs:
  :type: git
  :repo_location: '.'
```

Now, when this is done, you should be able to run PAC, the example below is where we extracted the idTestRepository to my home folder:

`docker run --rm -v /home/youruser/idReportTestRepository:/data praqma/pac from f9a66ca6d2e6`

If PAC is working, you should see the following on system out:

	[PAC] Applying task system none

and if you do an `ls -al` in your repostitory it should now look like this:

	-rwxrwxrwx 1 mads mads   759 Apr 12 10:24 default.html
	-rwxrwxrwx 1 mads mads   356 Apr 12 10:24 default.md
	-rwxrwxrwx 1 mads mads 21047 Apr 12 10:24 default.pdf
	-rw-rw-r-- 1 mads mads   608 Apr 12 10:23 default_settings.yml
	drwxrwxr-x 8 mads mads  4096 Apr 27  2015 .git
	-rwxrwxrwx 1 mads mads   489 Apr 12 10:24 ids.md
	-rw-rw-r-- 1 mads mads   340 Apr 27  2015 README.md

That's it. You've now succesfully created a changelog, automagically. As an alterntive you can run the script we provide with PAC (`demo_setup_docker.sh`) this script replays what was explained above. 

