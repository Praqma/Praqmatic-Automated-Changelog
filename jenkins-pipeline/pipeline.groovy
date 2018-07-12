multibranchPipelineJob("Praqmatic Automated Changelog") {
    factory {
        workflowBranchProjectFactory {
            scriptPath('jenkins-pipeline/Jenkinsfile')
        }
    }
    branchSources {
        git {
            credentialsId("github")
            remote("https://github.com/Praqma/Praqmatic-Automated-Changelog.git")
        }

        triggers {
            periodic(20)
        }
    }

    configure {
        it / 'sources' / 'data' / 'jenkins.branch.BranchSource' / 'source' / 'traits' << 'jenkins.plugins.git.traits.CloneOptionTrait' {
            extension(class: 'hudson.plugins.git.extensions.impl.CloneOption') {
                shallow(false)
                noTag(false)
                reference()
                depth(0)
                honorRefspec(false)
            }
        }
    }
}
