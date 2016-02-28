var request = require('request');

var Scales = {
    CHROMATIC: [0,1,2,3,4,5,6,7,8,9,10,11,12],
    MAJOR: [0,2,4,5,7,9,11,12]
}

var VOLTAGES = [0,5,10,15,20,25,30,35,40,45,50,55,60]
// #define octave2 120


function sendUdp(message) {
    var dgram = require("dgram");
    var message = new Buffer(message);
    var client = dgram.createSocket( "udp4" );

    var host = "192.168.240.1", port = 5555;
    client.send(message, 0, message.length, port, host);
}

function ArduinoApi(options) {
    var getVoltageForNote = function(note) {
        return VOLTAGES[note]
    }

    var sendNote = function(note) {
        var voltage = getVoltageForNote(note)
        console.log("playing " + voltage)
        sendUdp("" + voltage)
    }

    return {
        sendNote: sendNote
    }
}

var api = new ArduinoApi({ outputPin: 6 })
var cursor = 0;
console.log("i did it")

var playNextNote = function() {
    var scale = Scales.MAJOR
    var note = scale[cursor]
    cursor = (cursor + Math.floor(Math.random() * 6) - 2 + scale.length) % scale.length
    // cursor = Math.floor(Math.random() * scale.length)
    // cursor = ++cursor % scale.length//Math.floor(Math.random() * scale.length)
    api.sendNote(note)
    setTimeout(playNextNote, 200)
}

playNextNote()
