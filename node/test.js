var child = require('child_process').fork('test2.js');

var net = require('net');
var server = net.createServer();
server.on('connection', function(net){
    console.log('xxxx parent');
    net.end('parent')
})
 server.listen(8887, function(){
    console.log('server bound');
    child.send('sever', server)
})