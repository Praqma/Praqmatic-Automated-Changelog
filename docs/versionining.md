# Version reporting in PAC

We want PAC reports the version number, when asked. Currently that is the help switch or if no arguments are given.

The major-minor-patch version in PAC comes from the commited file called version.properties.

During the release process we write a file called version.stamp which contain one line:

$VER-$BUILD_NUMBER ($SHA)

e.g. 2.0.0-142 (ac6f04a)

The version.stamp is not commited back to the repository and created for each release process.

For the Docker image build as part of the release process this means the version.stamp file is part of the image.
For plain Ruby script used directly the version will be unknown. This is to some extend true, as we never know if the workspace is clean when running from a git repository.

Developer control the version.properties - nothing else.

When PAC print helps, it checks if version.stamp exists, and then prints it content integrated in the help and usage message.
If the file do not exists, it prints unknown version.
