pipeline {
  agent {
    docker {
      image 'golang'
      args '-v $WORKSPACE:/go/src/github.com/inhandnet/elements-cli'
    }
    
  }
  stages {
    stage('build') {
      steps {
        sh '''pwd
ls -al
go build github.com/inhandnet/elements-cli'''
      }
    }
    stage('stage') {
      steps {
        archiveArtifacts 'elements-cli'
      }
    }
  }
}