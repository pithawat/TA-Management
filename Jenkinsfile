pipeline {
    agent any

   environment {
        TEST_DB_COMPOSE = 'docker-compose.test.yml'
        APP_NAME = 'ta-management'
        DOCKER_IMAGE_TAG = "${APP_NAME}:${env.BUILD_NUMBER}"
        DOCKER_REGISTRY = 'ghcr.io/pithawat'
        FULL_IMAGE_NAME = '${DOCKER_REGISTRY}/${DOCKER_IMAGE_TAG}'
   }

   stages {
    stage('Checkout') {
        steps{
            checkout scm
            echo "Repository checked out successfully"
        }
    }



    stage('Configure Enviroment and Run Tests'){
        steps{
            withCredentials([
                file(credentialsId: 'DOT_ENV_FILE', variable: 'ENV_PATH')
            ]){
                sh ('''
                    mv "${ENV_PATH}" ./.env
                    ls -al ./.env
                ''')
            }

        }
    }

    stage('Build Image'){
        steps{
            script{
                echo "Building the application binary and testing image.."

                sh "docker build -t ${DOCKER_IMAGE_TAG} --target test-builder ."

                sh "docker tag ${DOCKER_IMAGE_TAG} ${FULL_IMAGE_NAME}"
            }
        }
    }

    stage('Integration Tests') {
        steps{
            script{
                echo "Starting Dockerized Integration Tests..."
                sh "docker compose -f ${TEST_DB_COMPOSE} down -v"
              // Capture exit code
                def testExitCode = sh(
                    script: "docker compose -f ${TEST_DB_COMPOSE} up --build --force-recreate --abort-on-container-exit --exit-code-from app_test",
                    returnStatus: true
                )
                
                // Check if tests failed
                if (testExitCode != 0) {
                    error("Integration tests failed with exit code: ${testExitCode}")
                }
                
                echo "✅ All integration tests passed!"

            }
        }
    }

    stage('Publish Image'){
        steps{
            script{
                withCredentials([usernamePassword(credentialsId: 'ghcr-creds', usernameVariable: 'GH_USER', passwordVariable: 'GH_PAT')]){
                sh """
                    echo "Logging into Github container Registry..."
                    echo \$GH_PAT | docker login ghcr.io -u \$GH_USER --password-stdin

                    echo "Tagging image ${DOCKER_IMAGE_TAG} to ${FULL_IMAGE_NAME}"
                    docker tag ${DOCKER_IMAGE_TAG} ${FULL_IMAGE_NAME}

                    echo "Pushing image to GHCR..."
                    docker push ${FULL_IMAGE_NAME}
                """
            }
            echo "Image successfully pushed to GHCR."
            }
        }
    }

    stage('Cleanup') {
            // Good practice: Ensure the temporary workspace is cleaned up to remove the secret file.
            steps {
                cleanWs()
            }
    }

    stage('Deploy to VM-TEST'){
        agent {label 'vm-test'}
        steps{
            echo "Deploying ${FULL_IMAGE_NAME} to Production Environment: ${env.NODE_NAME}"

            withCredentials([usernamePassword(credentialsId: 'ghcr-creds', usernameVariable: 'GH_USER', passwordVariable: 'GH_PAT')]){
                sh """
                    echo "Logging into Github container Registry..."
                    echo \$GH_PAT | docker login ghcr.io -u \$GH_USER --password-stdin
                    sudo usermod -aG docker $USER
                    newgrp docker
                    docker pull ${FULL_IMAGE_NAME}

                    docker stop $APP_NAME || true
                    docker rm $APP_NAME || true

                    docker run -d --name $APP_NAME -p 8084:8080 ${FULL_IMAGE_NAME}
                """
            }
         }
    }
    
    stage('Deploy to VM-PROD'){
        agent {label 'vm-prod'}
        steps{
            echo "Deploying ${FULL_IMAGE_NAME} to Production Environment: ${env.NODE_NAME}"

            withCredentials([usernamePassword(credentialsId: 'ghcr-creds', usernameVariable: 'GH_USER', passwordVariable: 'GH_PAT')]){
                sh """
                    echo "Logging into Github container Registry..."
                    echo \$GH_PAT | docker login ghcr.io -u \$GH_USER --password-stdin
                    docker pull ${FULL_IMAGE_NAME}

                    docker stop $APP_NAME || true
                    docker rm $APP_NAME || true

                   

                """
            }
            withCredentials([
                file(credentialsId: 'DOT_ENV_FILE', variable: 'ENV_PATH')
            ]){
                sh ('''
                    pwd
                    mv "${ENV_PATH}" ./.env
                    ls -al ./.env
                    docker run -d --name $APP_NAME -p 8084:8084 ${FULL_IMAGE_NAME}
                ''')
            }
        }
    }

    stage('start DB'){
        agent {label 'vm-db'}
        steps{
            withCredentials([
                file(credentialsId: 'DOT_ENV_FILE', variable: 'ENV_PATH')
            ]){
                sh ('''
                    pwd
                    mv "${ENV_PATH}" ./.env
                    ls -al ./.env
                ''')
            }
            script{
                echo "Run DB container"
                sh "docker compose up -d"
            }
        }
    }
}
}