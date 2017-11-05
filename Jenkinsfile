pipeline {
  agent any
  stages {
    stage('Test') {
      steps {
        sh 'GOROOT="/opt/go" PATH="/opt/go/bin:$PATH" make test'
      }
    }
    stage('Build') {
      steps {
        sh 'GOROOT="/opt/go" PATH="/opt/go/bin:$PATH" GOOS=linux GOARCH=amd64 make'
        sh 'GOROOT="/opt/go" PATH="/opt/go/bin:$PATH" GOOS=darwin GOARCH=amd64 make'
        sh 'GOROOT="/opt/go" PATH="/opt/go/bin:$PATH" GOOS=windows GOARCH=amd64 make'
      }
    }
    stage('Archive') {
      steps {
        archiveArtifacts 'bin/darwin_amd64/pressure'
        archiveArtifacts 'bin/pressure'
        archiveArtifacts 'bin/windows_amd64/pressure.exe'
      }
    }
  }
}