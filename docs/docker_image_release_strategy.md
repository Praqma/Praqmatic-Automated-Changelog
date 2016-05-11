# Relase strategy and design for PAC docker image

This file explain our different consideration related to releasing PAC as a Docker image.

## Background

Prior version 2.x of Praqmatic Automated Changelog the praqma/pac container was only containing a Ruby environment that enabled the users of PAC to avoid handling different complex aspects of Ruby installations.
The praqma/pac container was also released from the [docker-pac git repository](https://github.com/Praqma/docker-pac).

When PAC 2.x version was planned to be released we wanted to release a Docker image with the PAC tool as entry point, so it became easier to use.
We also wanted to use the same image as development environment for PAC, so developers would have a matching Ruby environmnet and could develop and chnage it together with PAC itself.

Considerations:

* Where should the Docker file be placed - seperate repository, like docker-pac or the PAC tool repository?
* If not the same repository how is Dockerfile and Ruby environment kept in sycn (e.g. the Bundle file)
* If PAC tool is seperated from the container we will need to copy in the PAC tool from another repository when building the Docker image

In favor of one repository - in which case it will be the Praqmatic Automed Changelog repo was:

* Avoid to manage dependencies across two projects
* Having a seperate repo (the old docker-pac repo) just to contain a Dockerfile and a readme

In favor of having the Docker file with the PAC tool and during the Ruby script development is to allow developers to easily test and verify that PAC can run inside the container by building the image locally with `docker build -t praqma/pac:snapshot .` 
It also makes sure that the docker image will never be out of sync with the gems PAC requires


Further in context of release we can use the same release job for PAC and image.

Some of the consideration can be read in more details and as discussion on issue [#9](https://github.com/Praqma/Praqmatic-Automated-Changelog/issues/9) and [#18](https://github.com/Praqma/Praqmatic-Automated-Changelog/issues/18)


## Current solution 

* Dockerfile is added the root of this repository
* The release job for the PAC tool, which basically add a git tag with a version number also releases the PAC container
* PAC container is released as praqma/pac with the same name as the old container, but the old have been deprecated updating readme etc. See section below.

### Docker hub autobuild vs. Jenkins builds


There are two obvious was to build and release the PAC image. One is using Docker hub autobuilds or the other one is to use our own Jenkins and some Jenkins docker plugins.

We did a short evaluation, trying to look at pro and cons. The following are repeated from [#18](https://github.com/Praqma/Praqmatic-Automated-Changelog/issues/18) to some extend.

**Docker hub autobuild**

* Pros
 * Uses Docker hub resources for building
 * Nice and neat application links, markdown highlighed documentation on Dockerhub (automatically)
 * Easy to set up tagging
* Cons
 * Needs setup to notify if build failed - else we may risk releasing a PAC tool version without a working image
 * Cannot easily be part of a pipeline, since the release should be a manual step...but this is done automatically. The autobuild only verifies that the image can be built but test is not part of this.
 * We'll need to move pac to new repository and replace the existing pac repo with the autobuild repository if we wish to maintain the shorthand praqma/pac notation for the image.

**Jenkins build**

* Pros
 * Very easy setup. We can borrow the publish step from existing DSL sources
 * Possible to use it as a part of our pipeline in a build step.
 * We can verify the image that will be built before publish
* Cons
 * No nice documention on dockerhub.
 * Requires machine resources


## Deprecating old docker-pac repository and praqma/pac containers

The older version of the praqma/pac image (tagged v3, v4, v18, v20) are only Ruby environments fitted to run older versions of PAC. In those versions you would have to mount in the PAC Ruby script yourself.
Those images do not work with never versions (2.x and forward).

The image was releases from another github repository (https://github.com/Praqma/docker-pac) which have now been deprecated (old versions still there if you check them out).

The old way of using this container is described here (for tag v20): https://github.com/Praqma/docker-pac/tree/f0d3c1300e0c03e86d310b2915be246ffade22a3
Check old releases here: https://github.com/Praqma/docker-pac/releases
