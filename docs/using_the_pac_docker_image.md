# Using the PAC docker image

We supply a PAC docker image, [`praqma/pac`](https://hub.docker.com/r/praqma/pac/), to easily run PAC and avoid any Ruby environment configuration.

## Try it!

 1. Pull the image: `docker pull praqma/pac` 
 2. Test-run the image: `docker run --rm praqm/pac` which when successful should output help and usage similar to this:

```
praqmatic automated changelog 

Usage:
  /usr/bin/pac (-d | --date) <to> [<from>] [--settings=<settings_file>] [--pattern=<rel_pattern>]     
  /usr/bin/pac (-s | --sha) <to> [<from>] [--settings=<settings_file>] [--pattern=<rel_pattern>]
  /usr/bin/pac (-t | --tag) <to> [<from>] [--settings=<settings_file>] [--pattern=<rel_pattern>]
  /usr/bin/pac -h|--help

Options:
  -h --help
             
    Show this screen.
    
  -d --date
             
    Use dates to select the changesets.
     
  -s --sha
              
    Use SHAs to select the changesets.      
  
  --settings=<settings_file> 
  
    Path to the settings file used. If nothing is specified default_settings.yml is used
    
  --pattern=<rel_pattern>
  
    Format that describes how your release tags look. This is used together with -t LATEST. We always check agains HEAD/TIP.
```

## Usage

With the PAC docker container, [the basic usage examples in the README becomes:](../README.md#usage)

		docker run --rm praqma/pac -h
    docker run --rm -v $(pwd):/data -v $(pwd):/pac-templates praqma/pac --settings=/data/pac_settings.yml -t Release-1.0
		docker run --rm -v $(pwd):/data -v $(pwd):/pac-templates praqma/pac --settings=/data/pac_settings.yml -t Release-1.0 Release-2.0
    docker run --rm -v $(pwd):/data -v $(pwd):/pac-templates praqma/pac --settings=/data/pac_settings.yml -d 2013-10-01

Try out a more elaborate example: [Try PAC with a test repo and PAC docker image](try_pac_with_test_repo_and_docker.md)

## Use it with a repository


 1. First write a file as outlined in the [README](../README.md#simple-configuration-file)
 2. Then create a template for the changelog report as outlined in the [README](../README.md#simple-template) 

Then we can try to use PAC. 

Assume we have our nice templates stored in our home directory under `~/pac-templates` and each project (repository) have its own PAC configuration file in the root of the repository and named `pac_settings.yml`.

Then when we run the docker container we mount in the git repository together with the templates:

```
docker run --rm -v $(pwd):/data -v ~/pac-templates:/pac-templates --settings=/data/pac_settings.yml praqma/pac -d 2013-01-01
```
Running this command will produce a report that takes all commits from 2013 up till `HEAD` of the repository.

The relevant PAC configuration file matching the above needs templates defined as follows if we assume you've created a template `~/pac-templates/my-template.md`:

```
  :templates:
    - { location: '/pac-templates/my-template.md', output: my-changelog.md }
```

That'll produce a `my-changelog.md` in the root of your your project. 