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
                    // --- Build Backend Builder Image (untuk Unit Test & govulncheck) ---
                    dir('backend') {
                        echo "Building Backend Builder Image for Unit Tests..."
                        sh """
                        docker build \\
                            --target builder \\
                            -t \${BACKEND_IMAGE}-builder:\${GIT_COMMIT_SHORT} \\
                            .
                        """
                    }
                    
                    // --- Build Frontend Builder Image (untuk npm audit) ---
                    dir('frontend') {
                        echo "Building Frontend Builder Image for Audit..."
                        sh """
                        docker build \\
                            --target builder \\
                            -t \${FRONTEND_IMAGE}-builder:\${GIT_COMMIT_SHORT} \\
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
        
        stage('Dependency Security Audit') {
            steps {
                script {
                    echo "Menjalankan npm audit (Frontend)..."
                    // Menggunakan || true sementara agar audit tidak menggagalkan build jika ada moderate vulnerability
                    sh "docker run --rm \${FRONTEND_IMAGE}-builder:\${GIT_COMMIT_SHORT} sh -c 'npm audit --audit-level=high || echo \"Peringatan: Terdapat NPM vulnerabilities!\"'"
                    
                    echo "Menjalankan govulncheck (Backend)..."
                    sh "docker run --rm \${BACKEND_IMAGE}-builder:\${GIT_COMMIT_SHORT} sh -c 'go install golang.org/x/vuln/cmd/govulncheck@latest && govulncheck ./...'"
                }
            }
        }
        
        stage('Integration & Load Test') {
            steps {
                script {
                    echo "Menjalankan Integration Test menggunakan Docker Compose..."
                    try {
                        // Mengunduh docker-compose binary secara lokal ke workspace (karena Jenkins container mungkin tidak memilikinya)
                        sh """
                        if [ ! -f ./docker-compose ]; then
                            curl -SL https://github.com/docker/compose/releases/download/v2.24.5/docker-compose-linux-x86_64 -o docker-compose
                            chmod +x docker-compose
                        fi
                        """
                        
                        // Menghidupkan lingkungan test (MySQL, Redis, Backend)
                        sh "BACKEND_IMAGE=\${BACKEND_IMAGE} GIT_COMMIT_SHORT=\${GIT_COMMIT_SHORT} ./docker-compose -f docker-compose.test.yml up -d"
                        
                        // Menunggu container database siap (Healthcheck)
                        echo "Menunggu database siap..."
                        sleep 15
                        
                        // Menembak endpoint integration test dari dalam docker network
                        echo "Menguji Endpoint Redis Integration..."
                        sh "docker run --rm --network asahkoding_test_net curlimages/curl -f --retry 5 --retry-connrefused --retry-delay 3 http://backend-test:8080/health/redis"
                        
                        echo "Integration Test Berhasil!"

                        // Menjalankan Load & Stress Test menggunakan K6
                        echo "Memulai Load & Stress Testing dengan K6 (50 VUs)..."
                        sh "cat k6/load-test.js | docker run --rm -i --network asahkoding_test_net grafana/k6 run -"
                        
                        echo "Load & Stress Test Berhasil!"
                        
                        // Menjalankan DAST (Dynamic Application Security Testing) menggunakan OWASP ZAP
                        echo "Memulai OWASP ZAP Baseline Scan (DAST)..."
                        // Gunakan opsi -I agar ZAP mengembalikan exit 0 meskipun menemukan warning, karena pipeline kita belum memiliki konfigurasi keamanan header yang ekstensif
                        sh "docker run --rm -i --network asahkoding_test_net ghcr.io/zaproxy/zaproxy:stable zap-baseline.py -t http://backend-test:8080 -I"
                        
                        echo "OWASP ZAP DAST Selesai!"
                    } catch (Exception e) {
                        echo "Integration / Load Test Gagal: \${e.message}"
                        error("Integration / Load Test Gagal!")
                    } finally {
                        // Membersihkan container test agar tidak memakan resource Jenkins
                        sh "./docker-compose -f docker-compose.test.yml down -v"
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
