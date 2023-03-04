test: 
	go build -o gossip ./broadcast/main.go && cd maelstrom/ && ./maelstrom test -w broadcast --bin ../gossip --node-count 1 --time-limit 20 --rate 10