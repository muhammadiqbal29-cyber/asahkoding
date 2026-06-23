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
                    // --- Build Backend Builder Image (untuk Unit Test) ---
                    dir('backend') {
                        echo "Building Backend Builder Image for Unit Tests..."
                        sh """
                        docker build \\
                            --target builder \\
                            -t \${BACKEND_IMAGE}-builder:\${GIT_COMMIT_SHORT} \\
                            .
                        """
                    }
                }
            }
        }
        
        stage('Unit Test & Coverage') {
            steps {
                script {
                    echo "Menjalankan Golang Unit Test & Coverage..."
                    // Menjalankan testing dari dalam image builder yang memiliki source code & tools Golang
                    sh """
                        docker run --rm \${BACKEND_IMAGE}-builder:\${GIT_COMMIT_SHORT} sh -c '
                            go test ./... -coverprofile=coverage.out
                            COVERAGE=\$(go tool cover -func=coverage.out | grep total | awk "{print \\\$3}" | tr -d "%")
                            echo "Current Coverage: \${COVERAGE}%"
                            awk -v cov="\$COVERAGE" "BEGIN { if (cov < 20.0) { exit 1 } }" || { echo "Coverage is below 20%! Failing build."; exit 1; }
                            echo "Coverage OK."
                        '
                    """
                }
            }
        }
        
        stage('Security Scan (Trivy)') {
            steps {
                script {
                    echo "Memindai Backend Image untuk Vulnerabilities..."
                    // Menjalankan Trivy dari Docker untuk mengecek image (hanya menampilkan HIGH dan CRITICAL)
                    sh "docker run --rm -v /var/run/docker.sock:/var/run/docker.sock aquasec/trivy image --severity HIGH,CRITICAL --no-progress \${BACKEND_IMAGE}:\${GIT_COMMIT_SHORT}"
                }
            }
        }
        
        stage('Integration Test') {
            steps {
                script {
                    echo "Menjalankan Integration Test menggunakan Docker Compose..."
                    try {
                        // Menghidupkan lingkungan test (MySQL, Redis, Backend)
                        sh "docker compose -f docker-compose.test.yml up -d"
                        
                        // Menunggu container database siap (Healthcheck)
                        echo "Menunggu database siap..."
                        sleep 15
                        
                        // Menembak endpoint integration test
                        echo "Menguji Endpoint Redis Integration..."
                        sh "curl -f --retry 5 --retry-connrefused --retry-delay 3 http://localhost:8081/health/redis"
                        
                        echo "Integration Test Berhasil!"
                    } catch (Exception e) {
                        echo "Integration Test Gagal: \${e.message}"
                        error("Integration Test Gagal!")
                    } finally {
                        // Membersihkan container test agar tidak memakan resource Jenkins
                        sh "docker compose -f docker-compose.test.yml down -v"
                    }
                }
            }
        }
        
        // Catatan: Tahapan selanjutnya (Push ke Registry, Deploy) 
        // akan disesuaikan saat Anda melanjutkan checklist berikutnya.
    }
    
    post {
        always {
            // Membersihkan workspace setelah build selesai agar tidak memakan ruang disk
            cleanWs()
        }
    }
}
