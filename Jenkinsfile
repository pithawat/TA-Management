pipeline {
   agent {
        docker {
            image 'golang:latest' // or your build image
            // This is the key option: mount the host socket
            args '-v /var/run/docker.sock:/var/run/docker.sock' 
        }
    }

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

                // This command runs the test runner and waits for its exit code.
                sh "docker compose -f ${TEST_DB_COMPOSE} up --build --force-recreate --abort-on-container-exit --exit-code-from app_test"

            }
        }
    }

    stage('Publish Image'){

        when{ expression {return currentBuild.result == "SUCCESS"}}
        steps{
            script{
                withCredentials([usernamePassword(credentialsId: 'ghcr-creds', usernameVariable: 'GH_USER', passwordVariable: 'GH_PAT')])
                sh """
                    echo "Logging into Github container Registry..."
                    echo \$GH_PAT | docker login ghcr.io -u \$GH_USER --password-stdin

                    echo "Tagging image ${IMAGE_NAME} to ${FULL_IMAGE_NAME}"
                    docker tag ${IMAGE_NAME} ${FULL_IMAGE_NAME}

                    echo "Pushing image to GHCR..."
                    docker push ${FULL_IMAGE_NAME}
                """
            }
            echo "Image successfully pushed to GHCR."
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
                sh '''
                    echo "Logging into Github container Registry..."
                    echo \$GH_PAT | docker login ghcr.io -u \$GH_USER --password-stdin
                    docker pull ${FULL_IMAGE_NAME}

                    docker stop $CONTAINER_NAME || true
                    docker rm $CONTAINER_NAME || true

                    docker run -d --name $CONTAINER_NAME -p 8080:8080 $REGISTRY/$IMAGE_NAME:$TAG
                '''
            }
         }
    }
    
    stage('Deploy to VM-PROD'){
        agent {label 'vm-prod'}
        steps{
            echo "Deploying ${FULL_IMAGE_NAME} to Production Environment: ${env.NODE_NAME}"

            withCredentials([usernamePassword(credentialsId: 'ghcr-creds', usernameVariable: 'GH_USER', passwordVariable: 'GH_PAT')]){
                sh '''
                    echo "Logging into Github container Registry..."
                    echo \$GH_PAT | docker login ghcr.io -u \$GH_USER --password-stdin
                    docker pull ${FULL_IMAGE_NAME}

                    docker stop $CONTAINER_NAME || true
                    docker rm $CONTAINER_NAME || true

                    docker run -d --name $CONTAINER_NAME -p 8080:8080 $REGISTRY/$IMAGE_NAME:$TAG

                    //run db prod
                    docker compose up -d
                '''
            }
        }
    }
}
}