pipeline {
    agent any
    tools {
        go 'go1.16'
    }
    environment {
        GO114MODULE = 'on'
        CGO_ENABLED = 0 
        GOPATH = "${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}"
        XDG_CACHE_HOME = "${GOPATH}"
    }
    stages {        
        stage('Pre Test') {
            steps {
                echo 'Installing dependencies'
                sh 'go version'
                sh 'go get -u golang.org/x/lint/golint'
            }
        }

        stage('Build') {
            steps {
                echo 'Build silpht'
                sh 'go build ./cmd/silpht'
            }
        }

        stage('Test') {
            steps {
                withEnv(["PATH+GO=${GOPATH}/bin"]){
                    echo 'Running vetting'
                    sh 'go vet ./cmd/silpht'
                    echo 'Running linting'
                    sh 'golint ./cmd/silpht'
                    echo 'Running test'
                    sh 'go test -v ./cmd/* ./internal/* ./pkg/*'
                }
            }
        }
        
    }
}
