var fs = require('fs');
var crypto = require('crypto');

function decrypt(aseKey, inputFile){
    var fileBody = fs.readFileSync(inputFile)
    var  decipher =  crypto.createDecipheriv("aes-128-cfb",new Buffer(aseKey) , fileBody.slice(0,16))
    var recv = decipher.update(fileBody.slice(16))
    
    fs.writeFileSync(inputFile + ".n.ts", recv)
}

decrypt("0123456789123456", "1.ts.aes")