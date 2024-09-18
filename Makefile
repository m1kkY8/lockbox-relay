run:
	@go run main.go

update:
	@ssh prod "cd /home/prod/irc && bash restart.sh"

push: 
	@docker build . -t m1kky8/irc:prod --push
	@ssh prod "cd /home/prod/irc && bash restart.sh"
