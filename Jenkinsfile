pipeline {
    agent none 
    stages {
        // --- STAGE 1: Test Environment ---
        stage('Docker Test on VM-TEST') {
            // This stage will ONLY run on an agent with the label 'vm-test'
            agent {
                label 'vm-test'
            }
            steps {
                script {
                    echo "Starting Docker validation on Test Environment: ${env.NODE_NAME}"
                    // Run the 'hello-world' container
                    sh 'sudo docker run --rm hello-world'
                    echo "Test environment check passed."
                }
            }
        }
        
        // --- STAGE 2: Production Environment ---
        stage('Docker Test on VM-PROD') {
            input {
                message "Proceed to run on Production node (vm-prod)?"
                ok "Yes, Deploy to Prod"
            }
            
            // This stage will ONLY run on an agent with the label 'vm-prod'
            agent {
                label 'vm-prod'
            }
            steps {
                script {
                    echo "Starting Docker validation on Production Environment: ${env.NODE_NAME}"
                    
                    // Run the 'hello-world' container
                    sh 'sudo docker run --rm hello-world'
                    
                    echo "Production environment check passed."
                }
            }
        }
    }
}