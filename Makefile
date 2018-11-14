build:
	GOOS=linux go build -v -o webhook_linux .

docker-img: build
	docker build -f Dockerfile.local -t eu.gcr.io/involuted-smile-221511/webhook .

docker-push: docker-img
	docker push eu.gcr.io/involuted-smile-221511/webhook

redeploy: docker-push
	kubectl delete -f deploy.yaml
	kubectl apply -f deploy.yaml
	kubectl delete -f validating-webh-conf.yaml
	kubectl create -f validating-webh-conf.yaml