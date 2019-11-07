#!groovy
def dtrRegistryCreds = [credentialsId: "ba1807ec-6876-4676-bd53-dc0d09147af4", url: "https://dtr.corp-us-east-1.aws.dckr.io"]

pipeline {
    agent {
        label "ubuntu-1604-aufs-stable"
    }
    stages {
        stage("build") {
            steps{
                sh 'cp .env.example .env'
                sh '''
                GIT_COMMIT=`git rev-parse --short HEAD`
                BUILD_TIME=`date +%FT%T%z`
                docker-compose -f docker-compose.prod.yml build --build-arg GIT_COMMIT=${GIT_COMMIT} --build-arg BUILD_TIME=${BUILD_TIME} --build-arg VERSION=${GIT_COMMIT} mayday client
                '''
            }
        }
        stage("push") {
            when {
                branch 'master'
            }
            steps {
                withDockerRegistry(dtrRegistryCreds) {
                    sh "docker-compose -f docker-compose.prod.yml push mayday client"
                }
            }
        }
    }
}
