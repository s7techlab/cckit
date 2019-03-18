.: generate

generate:
	@echo "schema"
	@protoc -I=./  -I=../../../vendor \
    --govalidators_out=./ \
	--go_out=./ ./*.proto