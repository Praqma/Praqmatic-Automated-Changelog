job('1_pretested-integration') {
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

    wrappers {
        pretestedIntegration("SQUASHED", "2.0", "origin")
    }
    publishers {
        pretestedIntegration()
    }

    publishers {
        downstream('2_test', 'SUCCESS')
    }
}


job('2_test') {
    scm {
        git('https://github.com/Praqma/Praqmatic-Automated-Changelog.git', '2.0')
    }
    steps {
        rake() {
            task('test')
            installation('(Default)')
        }
    }
    publishers {
        downstream('3_functional_test', 'SUCCESS')
    }
}

job('3_functional_test') {
    scm {
        git('https://github.com/Praqma/Praqmatic-Automated-Changelog.git', '2.0')
    }
    steps {
        rake() {
            task('functional_test')
            installation('(Default)')
        }
    }
}

job('4_release') {

    scm {
        git('https://github.com/Praqma/Praqmatic-Automated-Changelog.git', '2.0')
    }


    publishers {
        git {
            pushOnlyIfSuccess()
            tag('origin', '$VERSION') {
                message('')
                create()
            }
        }
    }

    wrappers {
        environmentVariables {
            propertiesFile('./version.properties')
            env('VERSION', '$ver-$BUILD_NUMBER')

        }

    }
}

