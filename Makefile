.PHONY: build encrypt decrypt clean

build:
	go build -o myapp main.go

encrypt:
	go run main.go --file /Users/shubhamkumar/Downloads/imp-data.txt --key "bHB&fuw2GD*I@Bsd" --out /Users/shubhamkumar/Downloads/imp-dataenc.txt --encrypt

decrypt:
	go run main.go --file /Users/shubhamkumar/Downloads/imp-dataenc.txt --key "bHB&fuw2GD*I@Bsd" --out /Users/shubhamkumar/Downloads/imp-datade.txt

clean:
	rm -f myapp

