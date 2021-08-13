// 외부모듈 포함
const express = require('express');
const app = express();

// 패브릭연결설정
const { FileSystemWallet, Gateway } = require('fabric-network');

const fs = require('fs');
const path = require('path');
const ccpPath = path.resolve(__dirname, 'connection.json');
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');

const ccp = JSON.parse(ccpJSON); // json 직렬화 go-> unmarshal

// 서버설정
const PORT = 8080;
const HOST = '0.0.0.0';

// use static file
app.use(express.static(path.join(__dirname, 'views')));

// configure app to use body-parser
app.use(express.json());
app.use(express.urlencoded({ extended: false }));

// / GET index.html 페이지 라우팅
app.get('/', (req, res)=>{
    res.sendFile(__dirname + 'views/index.html');
})

async function cc_call(fn_name, args, res){
    
    // 인증서 가져오는부분
    const walletPath = path.join(process.cwd(), 'wallet');
    const wallet = new FileSystemWallet(walletPath);

    const userExists = await wallet.exists('admin');
    if (!userExists) {
        console.log('An identity for the user "admin" does not exist in the wallet');
        console.log('Run the registerUser.js application before retrying');
        return;
    }

    // 게이트웨이에 연결하는 부분 fabric-network + admin 인증서 + ccp (연결정보)
    const gateway = new Gateway();
    await gateway.connect(ccp, { wallet, identity: 'admin', discovery: { enabled: false } });
    // 채널에 연결
    const network = await gateway.getNetwork('mychannel');
    // 체인코드 연결
    const contract = network.getContract('luckydraw');

    var result;
    

    if(fn_name == 'register')
    {
       //console.log(args.toString())
        result = await contract.submitTransaction('register', args[0], args[1], args[2], args[3]);
    }
     
        
    else if( fn_name == 'join')
        result = await contract.submitTransaction('join', args[0], args[1]);
    else if(fn_name == 'draw')
        result = await contract.submitTransaction('draw', args[0], args[1]);
    else if(fn_name == 'finalize')
        result = await contract.submitTransaction('finalize', args[0]);
    else if(fn_name == 'query')
    {
        result = await contract.evaluateTransaction("query",args);
        console.log(result.toString())
        const myobj = JSON.parse(result)
        res.status(200).json(myobj)
    }
    else if(fn_name == 'history')
        result = await contract.evaluateTransaction('history', args[0]);
    else
        result = 'not supported function'

    return result;
}
// REST API 라우팅
// /draw  POST  라우팅 -> 추첨이벤트 등록
//pid, pname, pmanager,pparam {"Args":["register","D101", "summer event", "MGR1", "3"]}
app.post('/draw', async(req, res)=>{
    
    // POST method 인경우 변수가 문서 body영역에 담겨서 전달
    const pid = req.body.pid;
    const mode = req.body.mode;

    // (TO DO) 오류체크 -> 각 변수가 주어진 형식에 맞게 전달되었는지?

    if (mode == 0)
    {
        const pname = req.body.pname;
        const pmanager = req.body.pmanager;
        const pparam = req.body.pparam;
        // (TO DO) 오류체크 -> 각 변수가 주어진 형식에 맞게 전달되었는지?
        
        console.log("add event: " + pid+pname);

        result = cc_call('register', [pid,pname,pmanager,pparam])
    }
    if (mode == 1)
    {
        // 1. mode=1   조인
        //mode=1, pid, uid {"Args":["join","D101", "P2004"]}
        const pname = req.body.uid;
        // (TO DO) 오류체크 -> 각 변수가 주어진 형식에 맞게 전달되었는지?
        
        console.log("join event: " + pid+uid);

        result = cc_call('join', {pid,uid})
    }
    if (mode == 2)
    {
        // 2. mode=2   추첨
        //mode=2, pid, uid {"Args":["draw","D101", "P2004"]}
        const pname = req.body.uid;
        // (TO DO) 오류체크 -> 각 변수가 주어진 형식에 맞게 전달되었는지?
        
        console.log("draw event: " + pid+uid);

        result = cc_call('draw', {pid,uid})
    }
    if (mode == 3)
    {
        // 3. mode=3   종료//
        //mode=3, pid {"Args":["finalize","D101"]}     

        // (TO DO) 오류체크 -> 각 변수가 주어진 형식에 맞게 전달되었는지?
        
        console.log("finalize event: " + pid+pname);

        result = cc_call('finalize', pid)
    }

    // (TO DO) 체인코드의 결과를 분석하여 해당하는 JSON데이터를 클라이언트에게 전송
    // (TO DO) 체인코드 tx잘 해결되었을때 -> status code:200
    // (TO DO) 체인코드 tx처리에 오류가 발생했을때 -> status code:400번대
    
    const myobj = {result: "success"}
    res.status(200).json(myobj) 
})

// /draw GET 라우팅 -> 추첨이벤트 조회
//pid {"Args":["query","D101"]}     
app.get('/draw', async(req, res)=>{
   
    // GET method 인 경우 변수가 url query영역에 담겨서 전달
    const pid = req.query.pid;
    // (TO DO) 오류체크 -> 각 변수가 주어진 형식에 맞게 전달되었는지?

    console.log("query event: " + pid);

    const temp = cc_call('query', pid, res)

    // (TO DO) 체인코드의 결과를 분석하여 해당하는 JSON데이터를 클라이언트에게 전송
    // (TO DO) 체인코드 tx잘 해결되었을때 -> status code:200
    // (TO DO) 체인코드 tx처리에 오류가 발생했을때 -> status code:400번대
    
    // const myobj = JSON.parse(result)
    // res.status(200).json("myobj")
})
// /draw/history GET 라우팅 -> 추첨이벤트 이력조회
//pid {"Args":["history","D101"]}  


// 서버시작
app.listen(PORT, HOST);
console.log(`Running on http://${HOST}:${PORT}`);