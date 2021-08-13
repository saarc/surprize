#!/bin/bash

# 체인코드 설치
docker exec cli peer chaincode install -n luckydraw -v 1.0 -p github.com/luckydraw

# 체인코드 설치 확인
docker exec cli peer chaincode list --installed

# 체인코드 배포
docker exec cli peer chaincode instantiate -n luckydraw -v 1.0 -C mychannel -c '{"Args":[]}' -P 'OR ("Org1MSP.member","Org2MSP.member","Org3MSP.member","Org4MSP.member","Org5MSP.member")'

sleep 3

#체인코드 배포 확인
docker exec cli peer chaincode list --instantiated -C mychannel

# 체인코드 테스트 invoke 
docker exec cli peer chaincode invoke -n luckydraw -C mychannel -c '{"Args":["register","D101", "summer event", "MGR1", "3"]}'
sleep 3

docker exec cli peer chaincode invoke -n luckydraw -C mychannel -c '{"Args":["join","D101", "P2004"]}'
sleep 3

docker exec cli peer chaincode invoke -n luckydraw -C mychannel -c '{"Args":["draw","D101", "P2004"]}'
sleep 3

docker exec cli peer chaincode invoke -n luckydraw -C mychannel -c '{"Args":["finalize","D101"]}'
sleep 3

# 체인코드 테스트 query
docker exec cli peer chaincode query -n luckydraw -C mychannel -c '{"Args":["query","D101"]}'

docker exec cli peer chaincode query -n luckydraw -C mychannel -c '{"Args":["history","D101"]}'
