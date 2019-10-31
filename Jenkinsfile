#!groovy
def dtrRegistryCreds = [credentialsId: "ba1807ec-6876-4676-bd53-dc0d09147af4", url: "https://dtr.corp-us-east-1.aws.dckr.io"]

pipeline {
    agent {
        label "ubuntu-1604-aufs-stable"
    }
    stages {
        stage("build") {
            steps{
                sh 'docker-compose -f docker-compose.prod.yml build mayday client'
            }
        }
        stage("Push") {
            steps {
                withDockerRegistry(dtrRegistryCreds) {
                    sh "docker-compose -f docker-compose.prod.yml push mayday client"
                }
            }
        }
    }
}
