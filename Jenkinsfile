pipeline {
  agent any
  stages {
    stage('Run Tests') {
      steps {
        sh 'GOROOT="/opt/go" PATH="/opt/go/bin:$PATH" make test'
      }
    }
    stage('Build Test Report') {
      steps {
        sh 'GOROOT="/opt/go" PATH="/opt/go/bin:$PATH" make test-verbose 2>&1 | go-junit-report > report.xml'
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
        archiveArtifacts 'report.xml'
        junit 'report.xml'
      }
    }
  }
}
