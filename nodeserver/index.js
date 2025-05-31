const express = require("express");

const app1 = express();
app1.use(express.json())
app1.use(express.urlencoded({extended:true}));
app1.get('/', (req,res)=>{
    res.header({"node-js-server1" : "true"});
    res.send("hello from nodejs server")
});

app1.listen(3000, ()=>{
    console.log("Server started on http://localhost:3000");
})

const app2 = express();
app2.use(express.json())
app2.use(express.urlencoded({extended:true}));
app2.get('/', (req,res)=>{
    res.header({"node-js-server2" : "true"});
    res.send("hello from nodejs server")
});

app2.listen(3001, ()=>{
    console.log("Server started on http://localhost:3001");
})

const app3 = express();
app3.use(express.json())
app3.use(express.urlencoded({extended:true}));
app3.get('/', (req,res)=>{
    res.header({"node-js-server3" : "true"});
    res.send("hello from nodejs server")
});

app3.listen(3002, ()=>{
    console.log("Server started on http://localhost:3002");
})

// const app4 = express();
// app4.use(express.json())
// app4.use(express.urlencoded({extended:true}));
// app4.get('/', (req,res)=>{
//     res.header({"node-js-server4" : "true"});
//     res.send("hello from nodejs server")
// });

// app4.listen(3003, ()=>{
//     console.log("server started on http://localhost:3003");
// })