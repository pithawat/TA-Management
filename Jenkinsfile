pipeline {
    agent any

    environment {
        TEST_DB_COMPOSE = 'docker-compose.test.yml'
        APP_NAME = 'ta-management'
        DOCKER_IMAGE_TAG = "${APP_NAME}:${env.BUILD_NUMBER}"
        DOCKER_REGISTRY = 'ghcr.io/pithawat'
        // แก้ไข: ใช้ Double Quotes เพื่อให้ตัวแปรทำงาน
        FULL_IMAGE_NAME = "${DOCKER_REGISTRY}/${DOCKER_IMAGE_TAG}"
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Build & Test Image') {
            steps {
                script {
                    // Build builder stage เพื่อเตรียม binary
                    sh "docker build -t ${DOCKER_IMAGE_TAG} --target test-builder ."
                }
            }
        }

        stage('Integration Tests') {
            steps {
                script {
                    sh "docker compose -f ${TEST_DB_COMPOSE} down -v"
                    def testExitCode = sh(
                        script: "docker compose -f ${TEST_DB_COMPOSE} up --build --force-recreate --abort-on-container-exit --exit-code-from app_test",
                        returnStatus: true
                    )
                    if (testExitCode != 0) {
                        error("Integration tests failed")
                    }
                }
            }
        }

        stage('Publish Image') {
            steps {
                script {
                    withCredentials([usernamePassword(credentialsId: 'ghcr-creds', usernameVariable: 'GH_USER', passwordVariable: 'GH_PAT')]) {
                        sh """
                            docker build -t ${FULL_IMAGE_NAME} --target final .
                            echo \$GH_PAT | docker login ghcr.io -u \$GH_USER --password-stdin
                            docker push ${FULL_IMAGE_NAME}
                        """
                    }
                }
            }
        }

        // ย้าย Stage DB ขึ้นมาก่อน Deploy App
        stage('Prepare DB') {
            agent { label 'vm-db' }
            steps {
                sh "docker run --rm -v /home/link/jenkins/workspace/TA-management:/ws alpine sh -c 'rm -rf /ws/init.sql'"
                checkout scm
                sh "docker compose down -v && docker compose up -d"
            }
        }

        stage('Deploy to VM-PROD') {
            agent { label 'vm-prod' }
            steps {
                withCredentials([
                    usernamePassword(credentialsId: 'ghcr-creds', usernameVariable: 'GH_USER', passwordVariable: 'GH_PAT'),
                    file(credentialsId: 'DOT_ENV_FILE', variable: 'ENV_PATH')
                ]) {
                    sh """
                        echo \$GH_PAT | docker login ghcr.io -u \$GH_USER --password-stdin
                        docker pull ${FULL_IMAGE_NAME}
                        docker stop ${APP_NAME} || true
                        docker rm ${APP_NAME} || true
                        
                        # รันโดยส่ง Environment Variables เข้าไปโดยตรง (ดีกว่าส่งไฟล์ path ที่ไม่มีอยู่จริง)
                        docker run -d \\
                            --name ${APP_NAME} \\
                            -p 8084:8084 \\
                            --env-file ${ENV_PATH} \\
                            ${FULL_IMAGE_NAME}
                    """
                }
            }
        }
    }
    
    post {
        always {
            cleanWs()
        }
    }
}