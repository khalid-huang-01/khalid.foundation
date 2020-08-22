package main

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"testing"
	"time"
)

var (
	group = "/group"
	acls = zk.WorldACL(zk.PermAll)
)

func AddMember(conn *zk.Conn) error {
	prefix := fmt.Sprintf("%s/member-", group)
	fmt.Println("prefix:" , prefix)
	path, err := conn.CreateProtectedEphemeralSequential(prefix, []byte{}, acls)
	if err != nil {
		return err
	}
	fmt.Println("Path: ", path)
	return nil
}

func TestGroup(t *testing.T) {
	conn, _, err := zk.Connect([]string{"127.0.0.1"}, time.Second) //*10)
	defer conn.Close()
	if err != nil {
		t.Fatal("ERROR: ", err)
	}

	_,_, ech, err := conn.ExistsW(group)
	if err != nil {
		t.Fatal("ERROR: ", err)
	}

	go watchMemberChange(ech)
	_, err = conn.Create("/group", []byte{}, 0, acls)
	if err != nil {
		t.Error("ERROR: ", err)
	}
	err = AddMember(conn)
	if err != nil {
		t.Fatal("ERROR: ", err)
	}

	err = AddMember(conn)
	if err != nil {
		t.Fatal("ERROR: ", err)
	}
	time.Sleep(20 * time.Second)
}

func watchMemberChange(ech <-chan zk.Event) {
	event := <-ech
	fmt.Println("*******************")
	fmt.Println("path:", event.Path)
	fmt.Println("type:", event.Type.String())
	fmt.Println("state:", event.State.String())
	fmt.Println("-------------------")
}