pipeline {
  agent any
  stages {
    stage('s1') {
      steps {
        sh 'echo \'Hello World!\''
      }
    }
    stage('s2') {
      parallel {
        stage('s2') {
          steps {
            echo 'Message Printing'
          }
        }
        stage('s2alt') {
          steps {
            isUnix()
          }
        }
      }
    }
    stage('s3') {
      steps {
        mail(subject: 'Lala', body: 'Lolo')
      }
    }
  }
}