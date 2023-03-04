package main

import (
	"testing"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func TestNode_HandleTopology(t *testing.T) {
	type fields struct {
		Node      *maelstrom.Node
		messages  []int
		neighbors []string
	}
	type args struct {
		msg maelstrom.Message
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test",
			fields: fields{
				Node:      maelstrom.NewNode(),
				messages:  []int{},
				neighbors: []string{},
			},
			args: args{
				msg: maelstrom.Message{
					Body: []byte(`{
						"type": "topology",
						"topology": {
						  "n1": ["n2", "n3"],
						  "n2": ["n1"],
						  "n3": ["n1"]
						}
					  }`),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Node{
				Node:      tt.fields.Node,
				messages:  tt.fields.messages,
				neighbors: tt.fields.neighbors,
			}
			if err := n.HandleTopology(tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("Node.HandleTopology() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
