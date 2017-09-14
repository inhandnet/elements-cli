pipeline {
  agent {
    docker {
      image 'golang'
      args '-v .:/go/src/github.com/inhandnet/elements-cli'
    }
    
  }
  stages {
    stage('build') {
      steps {
        parallel(
          "build": {
            sh 'go build .'
            
          },
          "stage": {
            archiveArtifacts 'elements-cli'
            
          }
        )
      }
    }
  }
}