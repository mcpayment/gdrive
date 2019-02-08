node (label: 'linux-builder'){
    
    stage ('build') {
        id = sh (returnStdout: true, script: 'id -u').trim()
        sh 'mkdir -p build'
            
        sh """docker run \
                --user ${id}:0 \
                -v '${pwd()}/build:/home/gdrive' \
                -e GOPATH="/home/gdrive" \
                -e GOOS=linux -e GOARCH=386 \
                --rm golang go get -v github.com/mcpayment/gdrive"""
    }
    
    stage('publish'){
        if (env.BRANCH_NAME == 'master'){
            println "Uploading binary to S3 repo..."
            withAWS(credentials: '24222bc0-90e5-4017-8168-983eaa3f32b5', region: 'ap-southeast-1') {
                s3Upload acl: 'PublicRead',
                        bucket: 'devopsdev16-misc-public-dependencies',
                        file: 'build/bin/linux_386/gdrive',
                        path: 'gdrive/gdrive'
            }
        }
    }
    
}