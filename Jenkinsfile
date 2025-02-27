import java.text.SimpleDateFormat;

pipeline {
  agent {
    kubernetes {
      yaml """
apiVersion: v1
kind: Pod
metadata:
  labels:
    jenkins: worker
spec:
  containers:
  - name: kaniko
    image: registry.cn-beijing.aliyuncs.com/shujiajia/executor:debug
    command:
    - sleep
    args:
    - 99999
    tty: true
    volumeMounts:
      - name: docker-config
        mountPath: /kaniko/.docker/
        readOnly: true
  - name: maven3
    image: maven:3-alpine
    command: ["/bin/sh","-c","sleep 100000"]
  volumes:
  - name: docker-config
    configMap:
      name: docker-config
"""
    }
  }
  environment {
    DATED_GIT_HASH = "${new SimpleDateFormat("yyMMddHHmmss").format(new Date())}${GIT_COMMIT.take(6)}"
  }
  stages {
    stage('Configure') {
      steps {
        echo "hello, starting"
      }
    }
    stage('Build with Kaniko') {
      steps {
        container('kaniko') {
              sh '/kaniko/executor -f `pwd`/Dockerfile -c `pwd` --cache=true \
              --destination=eff4858/httpserver:${DATED_GIT_HASH} \
                      --insecure \
                      --skip-tls-verify  \
                      -v=debug'
        }
      }
    }
    stage('build maven3'){
      steps {
        container('maven3') {
            sh 'ls ;`pwd`'
        }
      }
    }
  }
}
