# SVN support for PAC - statement of work

## Problem

Currently PAC only supports Git as a version control. We want to enable users of SVN to also automatically generate a changelog based on the same principle as in the Git implementation.  

## Solution

We'll integrate the new features in the existing application, all changes will be fully backwards compatible. 

## Implementation

The solution will be implemented using the Subversion command line tool. We'll make use of the tool's ability to output changes in an XML format for easier parsing. There are several similar tools out there which can be used as inspiration in the implementation. 

## Deliveries

We'll distribute the source code as-is on GitHub, and also build and publish a working docker image with the new svn features.  

## Work load

Estimated workload is about 37 hours, including test, release and publication.