pipeline {
    agent any

    parameters {
        booleanParam(name: 'RUN_HEAVY_TESTS', defaultValue: false, description: 'Jalankan Load Test (K6) dan Pen-Test (ZAP) (Hanya direkomendasikan untuk rilis/nightly)')
    }

    environment {
        // Variabel untuk Docker Hub / Registry Anda
        DOCKER_REGISTRY = 'iqbalmahad'
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
                        
                        echo "Building Cypress E2E Image..."
                        sh """
                        docker build \\
                            -f Dockerfile.e2e \\
                            -t \${FRONTEND_IMAGE}-e2e:\${GIT_COMMIT_SHORT} \\
                            .
                        """
                    }
                }
            }
        }
        
        stage('Parallel Tests & Scans') {
            parallel {
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
                            if (params.RUN_HEAVY_TESTS) {
                                echo "Memindai Backend Image untuk Vulnerabilities..."
                                // Menjalankan Trivy dari Docker untuk mengecek image (hanya menampilkan HIGH dan CRITICAL)
                                sh "docker run --rm -v /var/run/docker.sock:/var/run/docker.sock aquasec/trivy image --severity HIGH,CRITICAL --no-progress \${BACKEND_IMAGE}:\${GIT_COMMIT_SHORT}"
                            } else {
                                echo "Melewati Security Scan (RUN_HEAVY_TESTS=false)"
                            }
                        }
                    }
                }
                
                stage('Dependency Security Audit') {
                    steps {
                        script {
                            if (params.RUN_HEAVY_TESTS) {
                                echo "Menjalankan npm audit (Frontend)..."
                                // Menggunakan || true sementara agar audit tidak menggagalkan build jika ada moderate vulnerability
                                sh "docker run --rm \${FRONTEND_IMAGE}-builder:\${GIT_COMMIT_SHORT} sh -c 'npm audit --audit-level=high || echo \"Peringatan: Terdapat NPM vulnerabilities!\"'"
                                
                                echo "Menjalankan govulncheck (Backend)..."
                                sh "docker run --rm \${BACKEND_IMAGE}-builder:\${GIT_COMMIT_SHORT} sh -c 'govulncheck ./...'"
                            } else {
                                echo "Melewati Dependency Security Audit (RUN_HEAVY_TESTS=false)"
                            }
                        }
                    }
                }
            }
        }
        
        stage('Integration, Load & E2E Test') {
            steps {
                script {
                    if (params.RUN_HEAVY_TESTS) {
                        echo "Menjalankan Lingkungan Test menggunakan Docker Compose..."
                        try {
                            // Mengunduh docker-compose binary secara lokal ke workspace
                            sh """
                            if [ ! -f ./docker-compose ]; then
                                curl -SL https://github.com/docker/compose/releases/download/v2.24.5/docker-compose-linux-x86_64 -o docker-compose
                                chmod +x docker-compose
                            fi
                            """
                            
                            // Menghidupkan lingkungan test dan menunggu seluruh container siap (Healthcheck) via argumen --wait
                            sh "FRONTEND_IMAGE=\${FRONTEND_IMAGE} BACKEND_IMAGE=\${BACKEND_IMAGE} GIT_COMMIT_SHORT=\${GIT_COMMIT_SHORT} ./docker-compose -f docker-compose.test.yml up --wait -d"
                            
                            echo "Semua kontainer siap (Healthcheck Passed)!"
                            
                            // Menembak endpoint integration test dari dalam docker network
                            echo "Menguji Endpoint Redis Integration..."
                            sh "docker run --rm --network asahkoding_test_net curlimages/curl -f --retry 5 --retry-connrefused --retry-delay 3 http://backend-test:8080/health/redis"
                            
                            echo "Integration Test Berhasil!"

                            // Menjalankan Load & Stress Test menggunakan K6
                            echo "Memulai Load & Stress Testing dengan K6 (50 VUs)..."
                            sh "cat k6/load-test.js | docker run --memory=200m --rm -i --network asahkoding_test_net grafana/k6 run -"
                            
                            echo "Load & Stress Test Berhasil!"
                            
                            // Menjalankan DAST (Dynamic Application Security Testing) menggunakan OWASP ZAP
                            echo "Memulai OWASP ZAP Baseline Scan (DAST)..."
                            // Gunakan opsi -I agar ZAP mengembalikan exit 0 meskipun menemukan warning
                            sh "docker run --memory=1g --rm -i --network asahkoding_test_net ghcr.io/zaproxy/zaproxy:stable zap-baseline.py -t http://backend-test:8080 -I"
                            
                            echo "OWASP ZAP DAST Selesai!"

                            // Menjalankan E2E Testing menggunakan Cypress
                            echo "Memulai End-to-End (E2E) Testing menggunakan Cypress..."
                            // Cypress membutuhkan flag --e2e. Kita oper CYPRESS_BASE_URL agar menembak container frontend-test
                            sh "docker run --memory=1g --rm --network asahkoding_test_net -e CYPRESS_BASE_URL=http://frontend-test:3000 \${FRONTEND_IMAGE}-e2e:\${GIT_COMMIT_SHORT} cypress run --e2e"
                            
                            echo "Cypress E2E Test Berhasil!"
                        } catch (Exception e) {
                            echo "Integration / Load / E2E Test Gagal: \${e.message}"
                            error("Integration / Load / E2E Test Gagal!")
                        } finally {
                            // Membersihkan container test agar tidak memakan resource Jenkins
                            sh "./docker-compose -f docker-compose.test.yml down -v"
                        }
                    } else {
                        echo "Melewati tahap Integration, Load, DAST, dan E2E Test (RUN_HEAVY_TESTS=false)."
                    }
                }
            }
        }
        
        stage('Release & Semantic Versioning') {
            steps {
                script {
                    echo "Melakukan rilis versi secara otomatis..."
                    
                    // 1. Konfigurasi Identitas Git untuk komit otomatis
                    sh "git config user.name 'Jenkins CI'"
                    sh "git config user.email 'jenkins@asahkoding.internal'"

                    // 2. Instal pustaka semantic versioning
                    sh "npm install"

                    // 3. Eksekusi Semantic Versioning (Bumping versi, buat CHANGELOG, buat Tag Git)
                    sh "npm run release"

                    // Baca versi yang baru terbentuk
                    NEW_VERSION = sh(returnStdout: true, script: "node -p \"require('./package.json').version\"").trim()
                    echo "Versi Baru yang Dirilis: v\${NEW_VERSION}"

                    // 4. Otentikasi dan push kembali Changelog dan Tag ke GitHub
                    withCredentials([usernamePassword(credentialsId: 'github-cred', passwordVariable: 'GIT_TOKEN', usernameVariable: 'GIT_USER')]) {
                        // Atur URL remote untuk menggunakan Token Personal Access
                        sh "git remote set-url origin https://\${GIT_USER}:\${GIT_TOKEN}@github.com/muhammadiqbal29-cyber/asahkoding.git"
                        // Push ke branch main berikut dengan tag-tag nya
                        sh "git push --follow-tags origin main"
                    }

                    // 5. Memberi tag versi baru pada Docker Image
                    echo "Memberi tag pada Docker Image menjadi v\${NEW_VERSION}..."
                    sh "docker tag \${BACKEND_IMAGE}:\${GIT_COMMIT_SHORT} \${BACKEND_IMAGE}:v\${NEW_VERSION}"
                    sh "docker tag \${FRONTEND_IMAGE}:\${GIT_COMMIT_SHORT} \${FRONTEND_IMAGE}:v\${NEW_VERSION}"

                    // 6. Push Image ke Docker Hub menggunakan Credentials Docker
                    withCredentials([usernamePassword(credentialsId: 'dockerhub-cred', passwordVariable: 'DOCKER_PASSWORD', usernameVariable: 'DOCKER_USERNAME')]) {
                        echo "Autentikasi ke Docker Hub..."
                        sh "echo \$DOCKER_PASSWORD | docker login -u \$DOCKER_USERNAME --password-stdin"
                        
                        echo "Mengunggah Docker Image ke Docker Hub..."
                        sh "docker push \${BACKEND_IMAGE}:v\${NEW_VERSION}"
                        sh "docker push \${FRONTEND_IMAGE}:v\${NEW_VERSION}"
                        
                        // Menghapus kredensial sesi lokal setelah selesai
                        sh "docker logout"
                    }
                    
                    echo "🎉 RILIS BERHASIL! Versi v\${NEW_VERSION} sudah live di Github dan Docker Hub."
                }
            }
        }
    }
    
    post {
        always {
            // Membersihkan workspace setelah build selesai agar tidak memakan ruang disk
            cleanWs()
        }
    }
}
