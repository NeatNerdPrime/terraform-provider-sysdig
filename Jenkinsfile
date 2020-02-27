pipeline {
   agent {
      docker { image 'golang:1.14' }
   }

   environment {
      SYSDIG_MONITOR_API_TOKEN = credentials('tech-marketing-token-monitor-lab')
      SYSDIG_SECURE_API_TOKEN = credentials('tech-marketing-token-secure-lab')
   }

   stages {
      stage('Unit Tests') {
         steps {
            sh 'make test'
         }
      }
      stage('Acceptance Tests') {
          steps {
            sh 'make testacc'
          }
      }
   }
}

