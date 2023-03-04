test: 
	go build -o gossip ./broadcast/main.go && cd maelstrom/ && ./maelstrom test -w broadcast --bin ../gossip --node-count 5 --time-limit 20 --rate 10