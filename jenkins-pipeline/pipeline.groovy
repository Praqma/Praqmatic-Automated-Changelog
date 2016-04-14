REPOSITORY_URL = 'https://github.com/Praqma/Praqmatic-Automated-Changelog.git'
MAIN_BRANCH = 'master'
REMOTE_NAME = 'origin'
JOB_LABELS = 'jenkinsubuntu'
NUM_OF_BUILDS_TO_KEEP = 100

job('1_pretested-integration') {
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

    //First step: Do PAC execute? (No syntax errors?)
    steps {
        shell('ruby pac.rb')
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
        downstreamParameterized {
            trigger('2_test') {
                parameters {
                    gitRevision()
                }
            }
        }
    }
}


job('2_test') {
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

    steps {
        shell('rake test')
    }

    wrappers {
        buildName('${BUILD_NUMBER}#${GIT_REVISION,length=8}(${GIT_BRANCH})')
    }

    publishers {
        mailer('and@praqma.net', false, false)
        downstreamParameterized {
            trigger('3_functional_test') {
                parameters {
                    gitRevision()
                }
            }
        }
    }
}

job('3_functional_test') {
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

    steps {
        shell('rake functional_test')
    }

    wrappers {
        buildName('${BUILD_NUMBER}#${GIT_REVISION,length=8}(${GIT_BRANCH})')
    }

    publishers {
        buildPipelineTrigger('4_release') {
            parameters {
                gitRevision()
            }
        }
        mailer('and@praqma.net', false, false)
    }
}

job('4_release') {
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
