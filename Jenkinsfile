pipeline {
  agent any
  stages {
    stage('Test') {
      steps {
        sh 'GOROOT="/opt/go" PATH="/opt/go/bin:$PATH" make test'
      }
    }
    stage('Build for Linux') {
      parallel {
        stage('Build for Linux') {
          steps {
            sh 'GOROOT="/opt/go" PATH="/opt/go/bin:$PATH" GOOS=linux GOARCH=amd64 make'
          }
        }
        stage('Build for Mac OS X') {
          steps {
            sh 'GOROOT="/opt/go" PATH="/opt/go/bin:$PATH" GOOS=darwin GOARCH=amd64 make'
          }
        }
        stage('Build for Windows ') {
          steps {
            sh 'GOROOT="/opt/go" PATH="/opt/go/bin:$PATH" GOOS=windows GOARCH=amd64 make'
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
}