REPOSITORY_URL = 'https://github.com/Praqma/Praqmatic-Automated-Changelog.git'
MAIN_BRANCH = 'master'
REMOTE_NAME = 'origin'
JOB_LABELS = 'GiJeLiSlave'
NUM_OF_BUILDS_TO_KEEP = 100
GITHUB_PRAQMA_CREDENTIALS = '100247a2-70f4-4a4e-a9f6-266d139da9db'

PRETESTED_INTEGRATION_JOB_NAME = '1_pretested-integration_pac'
FUNCTIONAL_TEST_JOB_NAME = '2_functional_test_pac'
RELEASE_JOB_NAME = '3_release_pac'

DOCKER_REPO_NAME = 'praqma/pac'

job(PRETESTED_INTEGRATION_JOB_NAME) {
    logRotator {
        numToKeep(NUM_OF_BUILDS_TO_KEEP)
    }

    label(JOB_LABELS)

    properties {
        ownership {
            primaryOwnerId('and')
            coOwnerIds('man')
        }
    }

    authorization {
        permission('hudson.model.Item.Read', 'anonymous')
    }

    scm {
        git {
            remote {
                name(REMOTE_NAME)
                url(REPOSITORY_URL)
                credentials(GITHUB_PRAQMA_CREDENTIALS)
            }
            branch('origin/ready/**')

            extensions {
                wipeOutWorkspace()
                pruneBranches()
            }
        }
    }

    triggers {
        githubPush()
    }

    //First step: Can we build the docker container, and can we run unit tests?
    //This basically mimics developer behaviour
    steps {
        shell('''docker build -t praqma/pac:snapshot .
                 |docker run --entrypoint=/bin/sh --rm -v $(pwd):/data praqma/pac:snapshot -c rake test'''.stripMargin())
    }

    wrappers {
        buildName('${BUILD_NUMBER}#${GIT_REVISION,length=8}(${GIT_BRANCH})')
        pretestedIntegration("SQUASHED", MAIN_BRANCH, REMOTE_NAME)
    }

    publishers {
        pretestedIntegration()
        downstream('2_functional_test_pac', 'SUCCESS')
        mailer('and@praqma.net', false, false)
    }

}

job(FUNCTIONAL_TEST_JOB_NAME) {
    logRotator {
        numToKeep(NUM_OF_BUILDS_TO_KEEP)
    }

    label(JOB_LABELS)

    properties {
        ownership {
            primaryOwnerId('and')
            coOwnerIds('man')
        }
    }

    authorization {
        permission('hudson.model.Item.Read', 'anonymous')
    }

    scm {
        git {
            remote {
                name(REMOTE_NAME)
                url(REPOSITORY_URL)
            }
            branch(MAIN_BRANCH)
            extensions {}
        }
    }

    //This is a workaround until we get docker inside docker to run our functional test
    steps {
        shell('''#!/bin/bash
                 |. ~/.profile
                 |rake functional_test'''.stripMargin())
    }

    wrappers {
        buildName('${BUILD_NUMBER}#${GIT_REVISION,length=8}(${GIT_BRANCH})')
    }

    publishers {
        buildPipelineTrigger(RELEASE_JOB_NAME) {
            parameters {
                gitRevision()
            }
        }
        mailer('and@praqma.net', false, false)
    }
}

job(RELEASE_JOB_NAME) {
    label(JOB_LABELS)

    properties {
        ownership {
            primaryOwnerId('and')
            coOwnerIds('man')
        }
    }

    authorization {
        permission('hudson.model.Item.Read', 'anonymous')
    }

    scm {
        git {
            remote {
                name(REMOTE_NAME)
                url(REPOSITORY_URL)
                credentials(GITHUB_PRAQMA_CREDENTIALS)
            }
            branch(MAIN_BRANCH)
            extensions {}
        }
    }

    wrappers {
        buildName('${BUILD_NUMBER}#${GIT_REVISION,length=8}(${GIT_BRANCH})-ver=${ENV,var="VERSION"}')

        environmentVariables {
            propertiesFile('./version.sh')
            env('VERSION', '$ver-$BUILD_NUMBER')
        }
    }

    steps {
      dockerBuildAndPublish {
        repositoryName(DOCKER_REPO_NAME)
        tag('${VERSION}')
        registryCredentials('docker-hub-crendential')
        dockerHostURI('unix:///var/run/docker.sock')
        forcePull(false)
        createFingerprints(false)
        skipDecorate()
      }
    }    

    publishers {
        git {
            pushOnlyIfSuccess()
            tag(REMOTE_NAME, '$VERSION') {
                message('')
                create()
            }

            mailer('and@praqma.net', false, false)
        }
    }
}
