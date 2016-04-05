job('1_test') {
    scm {
        git('git://github.com/jgritman/aws-sdk-test.git', '2.0')
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

job('2_functional_test') {
    scm {
        git('git://github.com/jgritman/aws-sdk-test.git', '2.0')
    }
    steps {
        rake() {
            task('functional_test')
            installation('(Default)')
        }
    }
}

job('3_release') {
    parameters {
        stringParam('version')
        textParam('description')
    }

    scm {
        git('git://github.com/jgritman/aws-sdk-test.git', '2.0')
    }

    steps {
        shell('git tag -a ${version} -m "${description}"')
        shell('git push origin ${version}')
    }

}