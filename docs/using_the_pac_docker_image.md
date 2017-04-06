# Using the PAC docker image

We supply a PAC docker image, [`praqma/pac`](https://hub.docker.com/r/praqma/pac/), to easily run PAC and avoid any Ruby environment configuration.

## Try it!

 1. Pull the image: `docker pull praqma/pac` 
 2. Test-run the image: `docker run --rm praqm/pac` which when successful should output help and usage similar to this:

```
praqmatic automated changelog 

Usage:
  #{__FILE__} from <oldest-ref> [--to <newest-ref>] [options] [-v...] [-q...] [-c (<user> <password> <target>)]... 
  #{__FILE__} from-latest-tag <approximation> [--to <newest-ref>] [options] [-v...] [-q...] [-c (<user> <password> <target>)]...
  #{__FILE__} -h|--help

Options:
  -h --help  Show this screen.

  --from <oldest-ref>  Specify where to stop searching for commit. For git this takes anything that rev-parse accepts. Such as HEAD~3 / Git sha or tag name.

  --from-latest-tag  Looks for the newest commit that the tag with <approximation> points to.  
              
  --settings=<path>  Path to the settings file used. If nothing is specified default_settings.yml is used      

  --properties=<properties>  

    Allows you to pass in additional variables to the Liquid templates. Must be in JSON format. Namespaced under properties.* in 
    your Liquid templates. Referenced like so '{{properties.[your-variable]}}' in your templates.

    JSON keys and values should be wrapped in quotation marks '"' like so: --properties='{ "title":"PAC Changelog" }'      

  -v  More verbose output. Can be repeated to increase output verbosity or to cancel out -q

  -q  Less verbose output. Can be repeated for more silence or to cancel out -v

  -c  Override username and password. Example: `-c my_user my_password jira`. This will set username and password for task system jira.

```

## Usage

With the PAC docker container, [the basic usage examples in the README becomes:](../README.md#usage)

		docker run --rm praqma/pac -h
    docker run --rm -v $(pwd):/data -v $(pwd):/pac-templates praqma/pac --settings=/data/pac_settings.yml from Release-1.0
		docker run --rm -v $(pwd):/data -v $(pwd):/pac-templates praqma/pac --settings=/data/pac_settings.yml from Release-1.0 --to Release-2.0

Try out a more elaborate example: [Try PAC with a test repo and PAC docker image](try_pac_with_test_repo_and_docker.md)

## Use it with a repository


 1. First write a file as outlined in the [README](../README.md#simple-configuration-file)
 2. Then create a template for the changelog report as outlined in the [README](../README.md#simple-template) 

Then we can try to use PAC. 

Assume we have our nice templates stored in our home directory under `~/pac-templates` and each project (repository) have its own PAC configuration file in the root of the repository and named `pac_settings.yml`.

Then when we run the docker container we mount in the git repository together with the templates:

```
docker run --rm -v $(pwd):/data -v ~/pac-templates:/pac-templates --settings=/data/pac_settings.yml praqma/pac from HEAD~3
```
Running this command will produce a report that takes all commits from 3 commits back up till `HEAD` of the repository.

The relevant PAC configuration file matching the above needs templates defined as follows if we assume you've created a template `~/pac-templates/my-template.md`:

```
  :templates:
    - { location: '/pac-templates/my-template.md', output: my-changelog.md }
```

That'll produce a `my-changelog.md` in the root of your your project. 