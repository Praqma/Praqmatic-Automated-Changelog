job('1_pretested-integration') {
    logRotator {
        numToKeep(25)
    }

    scm {
        git {
            remote {
                name('origin')
                url('https://github.com/Praqma/Praqmatic-Automated-Changelog.git')
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

    wrappers {
        pretestedIntegration("SQUASHED", "master", "origin")
    }

    publishers {
        pretestedIntegration()
    }

    publishers {
        downstream('2_test', 'SUCCESS')
        mailer('and@praqma.net bue@praqma.net', false, false)
    }
}


job('2_test') {
    logRotator {
        numToKeep(25)
    }

    scm {
        git('https://github.com/Praqma/Praqmatic-Automated-Changelog.git', 'master')
    }

    steps {
        rake() {
            task('test')
            installation('(Default)')
        }
    }

    publishers {
        downstream('3_functional_test', 'SUCCESS')
        mailer('and@praqma.net bue@praqma.net', false, false)
    }
}

job('3_functional_test') {
    logRotator {
        numToKeep(25)
    }

    scm {
        git('https://github.com/Praqma/Praqmatic-Automated-Changelog.git', 'master')
    }

    steps {
        rake() {
            task('functional_test')
            installation('(Default)')
        }
    }

    publishers {
        buildPipelineTrigger('4_release')
        mailer('and@praqma.net bue@praqma.net', false, false)
    }
}

job('4_release') {
    scm {
        git {
            remote {
                name('origin')
                url('https://github.com/Praqma/Praqmatic-Automated-Changelog.git')
            }
            branch('master')
            extensions {}
        }
    }
    publishers {
        git {
            pushOnlyIfSuccess()
            tag('origin', '$VERSION') {
                message('')
                create()
            }

            mailer('and@praqma.net bue@praqma.net', false, false)
        }

        wrappers {
            environmentVariables {
                propertiesFile('./version.properties')
                env('VERSION', '$ver-$BUILD_NUMBER')
            }
        }
    }
}

