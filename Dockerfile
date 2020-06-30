FROM golang:1.12.17 as builder

RUN mkdir /app 
ADD . /app/
WORKDIR /app 
RUN go build -o faasbenchmark main.go
RUN go build -o faasbenchmark-tui tui.go

FROM node:13.8.0-stretch

RUN apt-get update && apt-get install -y ca-certificates curl apt-transport-https lsb-release gnupg wget software-properties-common gcc zip unzip python3 python3-pip vim

# Install kubectl and aws-iam-authenticator
RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl && \
	chmod +x ./kubectl && \
	mv ./kubectl /usr/local/bin && \
	curl -o aws-iam-authenticator https://amazon-eks.s3.us-west-2.amazonaws.com/1.15.10/2020-02-22/bin/linux/amd64/aws-iam-authenticator && \
	chmod +x ./aws-iam-authenticator && \
	mv ./aws-iam-authenticator /usr/local/bin

# Install faas-cli
RUN curl -sSL https://cli.openfaas.com | sh

# Install openwhisk-cli
RUN curl -LO https://github.com/apache/openwhisk-cli/releases/download/1.0.0/OpenWhisk_CLI-1.0.0-linux-amd64.tgz && \
	tar -xz -f OpenWhisk_CLI-1.0.0-linux-amd64.tgz && \
	rm OpenWhisk_CLI-1.0.0-linux-amd64.tgz && \
	mv ./wsk /usr/local/bin && \
	echo "APIGW_ACCESS_TOKEN=token" >> ~/.wskprops

# Install kubeless cli
RUN curl -LO https://github.com/kubeless/kubeless/releases/download/v1.0.6/kubeless_linux-amd64.zip && \
	unzip kubeless_linux-amd64.zip && \
	mv bundles/kubeless_linux-amd64/kubeless /usr/local/bin && \
	rm -rf bundles/ && rm kubeless_linux-amd64.zip

# Install fission cli
RUN curl -Lo fission https://github.com/fission/fission/releases/download/1.8.0/fission-cli-linux && \
	chmod +x fission && \
	mv fission /usr/local/bin/

# add azure cli repo
RUN curl -sL https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor | tee /etc/apt/trusted.gpg.d/microsoft.asc.gpg > /dev/null && \
	AZ_REPO=$(lsb_release -cs) && \
	echo "deb [arch=amd64] https://packages.microsoft.com/repos/azure-cli/ $AZ_REPO main" | tee /etc/apt/sources.list.d/azure-cli.list

# add dotnet repo
RUN wget -qO- https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor > microsoft.asc.gpg && \
	mv microsoft.asc.gpg /etc/apt/trusted.gpg.d/ && \
	wget -q https://packages.microsoft.com/config/debian/9/prod.list && \
	mv prod.list /etc/apt/sources.list.d/microsoft-prod.list && \
	chown root:root /etc/apt/trusted.gpg.d/microsoft.asc.gpg && \
	chown root:root /etc/apt/sources.list.d/microsoft-prod.list

RUN npm install -g serverless azure-functions-core-tools@3 --unsafe-perm=true --allow-root
RUN apt-get update && apt-get install azure-cli openjdk-8-jdk maven dotnet-sdk-3.1 libsecret-1-dev -y --fix-missing
RUN mkdir /app

COPY --from=builder /app/ /app
WORKDIR /app
RUN npm install

CMD ./faasbenchmark





