process.on('message', function(m, server){
    if(m === 'sever'){
        server.on('connection', function(net){
            console.log('xxxx child');
            net.end('chi233ld1sdas');
        })
    }
})