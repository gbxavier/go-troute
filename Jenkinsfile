#!/usr/bin/env groovy

def imageName = "sandbox"
def registry = "gbxavier"
def tag = "go-traceroute"

node('master'){

    stage ('checkout'){
        checkout scm
    }

    stage ('Install Dependencies') {
        sh "go get -d -v ./..."
    }

    stage ('Unit Tests'){
        sh "go test"
    }

    stage ('Docker Image Build'){
        sh "docker build -t $registry/$imageName:$tag ."
    }

    stage ('Publish'){
        sh "docker push $registry/$imageName:$tag"
    }

}