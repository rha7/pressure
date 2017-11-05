pipeline {
  agent any
  stages {
    stage('Test') {
      steps {
        sh 'make test'
      }
    }
    stage('Build for Linux') {
      parallel {
        stage('Build for Linux') {
          steps {
            sh 'GOOS=linux GOARCH=amd64 make'
          }
        }
        stage('Build for Mac OS X') {
          steps {
            sh 'GOOS=darwin GOARCH=amd64 make'
          }
        }
        stage('Build for Windows ') {
          steps {
            sh 'GOOS=windows GOARCH=amd64 make'
          }
        }
      }
    }
    stage('Timestamps') {
      steps {
        timestamps()
      }
    }
  }
  environment {
    GOROOT = '/opt/go1.8.3'
  }
}