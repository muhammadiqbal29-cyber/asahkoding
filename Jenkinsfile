pipeline {
    agent any

    environment {
        // Variabel untuk Docker Hub / Registry Anda
        DOCKER_REGISTRY = 'iqbalmahad4039'
        BACKEND_IMAGE = "${DOCKER_REGISTRY}/leetcode-backend"
        FRONTEND_IMAGE = "${DOCKER_REGISTRY}/leetcode-frontend"
        
        // Membaca versi commit Git pendek untuk tagging
        GIT_COMMIT_SHORT = sh(returnStdout: true, script: 'git rev-parse --short HEAD').trim()
    }

    stages {
        stage('Checkout') {
            steps {
                // Jenkins otomatis mengambil dari konfigurasi SCM Webhook
                checkout scm
            }
        }
        
        stage('Build Docker Images') {
            steps {
                script {
                    // --- Build Backend Image ---
                    dir('backend') {
                        echo "Building Backend Image: ${BACKEND_IMAGE}:${GIT_COMMIT_SHORT}"
                        // Memanfaatkan Docker Cache Layer (--cache-from)
                        sh """
                        docker build \
                            --cache-from ${BACKEND_IMAGE}:latest \
                            -t ${BACKEND_IMAGE}:${GIT_COMMIT_SHORT} \
                            -t ${BACKEND_IMAGE}:latest \
                            .
                        """
                    }
                    
                    // --- Build Frontend Image ---
                    dir('frontend') {
                        echo "Building Frontend Image: ${FRONTEND_IMAGE}:${GIT_COMMIT_SHORT}"
                        // Memanfaatkan Docker Cache Layer (--cache-from)
                        sh """
                        docker build \
                            --cache-from ${FRONTEND_IMAGE}:latest \
                            -t ${FRONTEND_IMAGE}:${GIT_COMMIT_SHORT} \
                            -t ${FRONTEND_IMAGE}:latest \
                            .
                        """
                    }
                }
            }
        }
        
        // Catatan: Tahapan selanjutnya (Test, Push ke Registry, Deploy) 
        // akan disesuaikan saat Anda melanjutkan checklist berikutnya.
    }
    
    post {
        always {
            // Membersihkan workspace setelah build selesai agar tidak memakan ruang disk
            cleanWs()
        }
    }
}
