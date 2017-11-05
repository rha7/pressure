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
        sh '''
            GOROOT="/opt/go" PATH="/opt/go/bin:$PATH" GOOS=linux GOARCH=amd64 make && \\
            GOROOT="/opt/go" PATH="/opt/go/bin:$PATH" GOOS=darwin GOARCH=amd64 make && \\
            GOROOT="/opt/go" PATH="/opt/go/bin:$PATH" GOOS=windows GOARCH=amd64 make && \\
            mv bin/pressure bin/pressure_linux_amd64 && \\
            mv bin/darwin_amd64/pressure bin/pressure_darwin_amd64 && \\
            mv bin/windows_amd64/pressure.exe bin/pressure_windows_amd64.exe
        '''
      }
    }
    post {
      always {
        archiveArtifacts 'bin/pressure_linux_amd64'
        archiveArtifacts 'bin/pressure_darwin_amd64'
        archiveArtifacts 'bin/pressure_windows_amd64.exe'
        archiveArtifacts 'report.xml'
        junit 'report.xml'
      }
    }
  }
}
