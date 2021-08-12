package main

// 외부모듈
import (
	"fmt"
	"encoding/json"
	"strconv"
	"time"
	"bytes"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// SmartContract 클래스 정의
type SmartContract struct {
}

// Draw 구조체 정의
type Draw struct {
	Pid				string		`json:"pid"`
	Pname			string		`json:"pname"`
	Pmanager		string		`json:"pmanager"`
	Pparam			int 		`json:"pparam`
	Pstate			string		`json:"pstate"` // registered, joining, draw, finalized
	Participants	[]string 	`json:"participants"`
	Winners			[]string	`json:"winners"`
}

// Init 함수 구현
func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) peer.Response {
	// 초기 시스템의 상태 설정
	// instantiate 권한 ca에 등록된 배포한사람의 role관리
	return shim.Success(nil)
}

// Invoke 함수 구현
// peer chaincode invoke -n luckydraw -C mychannel -c '{"Args":["finalize","D101"]}'
func (s *SmartContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {

	fn, args := stub.GetFunctionAndParameters()
	
	if fn == "register" {			// 블록생성
		return s.register(stub, args)
	} else if fn == "join" {			// 블록생성
		return s.join(stub, args)
	} else if fn == "draw" {			// 블록생성
		return s.draw(stub, args)
	} else if fn == "finalize" {		// 블록생성
		return s.finalize(stub, args)
	} else if fn == "query" {			// ws 조회
		return s.query(stub, args)
	} else if fn == "history" {		// 블록조회
		return s.history(stub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}
// register //  Creation R U D
// peer chaincode invoke -n luckydraw -C mychannel -c '{"Args":["register","D101", "summer event", "MGR1", "3"]}'
func (s *SmartContract) register(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of argument. Expecting 4")
	}

	// (TO DO) 이미 등록된 EVENT ID인가? 


	Nwinners, _ := strconv.Atoi(args[3])
	var draw = Draw{Pid: args[0], Pname: args[1], Pmanager: args[2], Pparam: Nwinners, Pstate:"registered"}

	drawAsBytes, _ := json.Marshal(draw)
	stub.PutState(args[0], drawAsBytes)

	return shim.Success(nil)
}
// join //  C R Update D
// peer chaincode invoke -n luckydraw -C mychannel -c '{"Args":["join","D101", "P2004"]}'
func (s *SmartContract) join(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of argument. Expecting 2")
	}
	
	// World State 조회
	drawAsBytes, _ := stub.GetState(args[0])
	if drawAsBytes == nil {
		return shim.Error("Requested draw id is missing")
	}
	// 객체화 (JSON -> 구조체)
	draw := Draw{}
	json.Unmarshal(drawAsBytes, &draw)

	// (TO DO) current state 가 registed or joining인가?

	// 수정 -> 참여자 추가
	draw.Participants = append(draw.Participants, args[1])
	draw.Pstate = "joining"

	// 직렬화 (구조체 -> JSON)
	drawAsBytes, _ = json.Marshal(draw)

	// World State 업데이트
	stub.PutState(args[0], drawAsBytes)

	return shim.Success(nil)
}
// draw //  C R Update D
// peer chaincode invoke -n luckydraw -C mychannel -c '{"Args":["draw","D101", "P2004"]}'
func (s *SmartContract) draw(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of argument. Expecting 2")
	}
	
	// World State 조회
	drawAsBytes, _ := stub.GetState(args[0])
	if drawAsBytes == nil {
		return shim.Error("Requested draw id is missing")
	}
	// 객체화 (JSON -> 구조체)
	draw := Draw{}
	json.Unmarshal(drawAsBytes, &draw)

	// (TO DO) 해당 위너가 Participant에 들어있나?, 이미 Winner인가?

	// (TO DO) 담청자의 수가 Pparam의 수보다 작은가? 3 -> 3 등록되어있으면 불가
	
	// (TO DO) current state = joining or draw 인가?

	// 수정 -> 참여자 추가
	draw.Winners = append(draw.Winners, args[1])
	draw.Pstate = "draw"

	// 직렬화 (구조체 -> JSON)
	drawAsBytes, _ = json.Marshal(draw)

	// World State 업데이트
	stub.PutState(args[0], drawAsBytes)

	return shim.Success(nil)
}
// finalize //  C R Update D
// peer chaincode invoke -n luckydraw -C mychannel -c '{"Args":["finalize","D101"]}'
func (s *SmartContract) finalize(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of argument. Expecting 2")
	}
	
	// World State 조회
	drawAsBytes, _ := stub.GetState(args[0])
	if drawAsBytes == nil {
		return shim.Error("Requested draw id is missing")
	}
	// 객체화 (JSON -> 구조체)
	draw := Draw{}
	json.Unmarshal(drawAsBytes, &draw)

	// (TO DO) current state = draw 인가?

	// 수정 -> 참여자 추가
	draw.Pstate = "finalized"

	// 직렬화 (구조체 -> JSON)
	drawAsBytes, _ = json.Marshal(draw)

	// World State 업데이트
	stub.PutState(args[0], drawAsBytes)

	return shim.Success(nil)
}

// query //  C Rear/Retrieve U D
// peer chaincode query -n luckydraw -C mychannel -c '{"Args":["query","D101"]}'
func (s *SmartContract) query(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	
	if len(args) != 1 {
		return shim.Error("Incorrect number of argument. Expecting 2")
	}
	
	// World State 조회
	drawAsBytes, _ := stub.GetState(args[0])
	if drawAsBytes == nil {
		return shim.Error("Requested draw id is missing")
	}
	return shim.Success(drawAsBytes)
}
// history
// peer chaincode query -n luckydraw -C mychannel -c '{"Args":["history","D101"]}'
func (t *SmartContract) history(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	drawid := args[0]

	fmt.Printf("- start history: %s\n", drawid)

	resultsIterator, err := stub.GetHistoryForKey(drawid)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the marble
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON marble)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- draw history returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

// main
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}