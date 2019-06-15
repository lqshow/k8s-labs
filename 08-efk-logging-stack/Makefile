.PHONY: deploy \
		deploy-elasticsearch \
		deploy-fluentd \
		deploy-kibana

deploy-elasticsearch:
	@kubectl apply -f es-statefulset.yaml
	@kubectl apply -f es-service.yaml

deploy-fluentd:
	@kubectl apply -f fluentd-es-configmap.yaml
	@kubectl apply -f fluentd-es-ds.yaml

deploy-kibana:
	@kubectl apply -f kibana-deployment.yaml
	@kubectl apply -f kibana-service.yaml

deploy: deploy-elasticsearch deploy-fluentd deploy-kibana