REPOSITORY_URL = 'https://github.com/Praqma/Praqmatic-Automated-Changelog.git'
MAIN_BRANCH = 'master'
REMOTE_NAME = 'origin'
JOB_LABELS = 'GiJeLiSlave'
NUM_OF_BUILDS_TO_KEEP = 100
def releasePraqmaCredentials = '100247a2-70f4-4a4e-a9f6-266d139da9db'

job('1_pretested-integration_pac') {
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
                credentials(releasePraqmaCredentials)
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
        shell('''
docker build -t praqma/pac:snapshot .             
docker run --entrypoint=/bin/sh --rm -v $(pwd):/data praqma/pac:snapshot -c rake test
''')
    }


    wrappers {
        buildName('${BUILD_NUMBER}#${GIT_REVISION,length=8}(${GIT_BRANCH})')
        pretestedIntegration("SQUASHED", MAIN_BRANCH, REMOTE_NAME)
    }

    publishers {
        pretestedIntegration()
    }

    publishers {
        mailer('and@praqma.net', false, false)
    }
}

job('2_functional_test_pac') {
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

    triggers {
        githubPush()
    }    

    //Since we do 'docker stuff' using rake...i don't know how tests would react if we start running docker in docker
    //TODO: This should be done differently. Since it requires manual configuration of a slave
    steps {
        shell('rake functional_test')
    }

    wrappers {
        buildName('${BUILD_NUMBER}#${GIT_REVISION,length=8}(${GIT_BRANCH})')
    }

    publishers {
        buildPipelineTrigger('3_release') {
            parameters {
                gitRevision()
            }
        }
        mailer('and@praqma.net', false, false)
    }
}

job('3_release_pac') {
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
                credentials(releasePraqmaCredentials)
            }
            branch(MAIN_BRANCH)
            extensions {}
        }
    }

    wrappers {
        buildName('${BUILD_NUMBER}#${GIT_REVISION,length=8}(${GIT_BRANCH})-ver=$VERSION')

        environmentVariables {
            propertiesFile('./version.sh')
            env('VERSION', '$ver-$BUILD_NUMBER')
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
