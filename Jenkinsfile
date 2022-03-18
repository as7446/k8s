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
          sh '/kaniko/executor -f `pwd`/golang/httpserver/Dockerfile -c `pwd`/golang/httpserver --cache=true \
          --destination=eff4858/httpserver:${DATED_GIT_HASH} \
                  --insecure \
                  --skip-tls-verify  \
                  -v=debug'
        }
      }
    }  
  }
}
