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
        sh 'go build github.com/inhandnet/elements-cli
      }
    }
    stage('stash') {
      steps {
        archiveArtifacts 'elements-cli'
      }
    }
  }
}