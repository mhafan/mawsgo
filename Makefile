lib:
	go build

push:
	git add *
	git commit -m 'Makefile push'
	git push

	