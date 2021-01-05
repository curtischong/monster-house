create-bucket:
	aws s3api --endpoint-url=http://localhost:4566 create-bucket --acl public-read --bucket monsterhouse
empty-bucket:
	aws --endpoint-url=http://localhost:4566 s3 rm s3://monsterhouse --recursive
list-buckets:
	aws --endpoint-url=http://localhost:4566 s3 ls
list-bucket-objects:
	aws --endpoint-url=http://localhost:4566 s3 ls monsterhouse
clean:
	# This deletes all of the data from the localstack environment. You will have to restart
	# the docker container to continue after running this command
	rm -rf ./data/data
run:
	go run main.go