var net = require('http');
var server = net.createServer((req, res)=>{
    res.end('welcome')
});
server.listen(8888, function(){
    console.log('server bound')
})