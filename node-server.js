var request = require('request'),
    _ = require('underscore'),
    improvisr = require("./improvisr.js");

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
        sendUdp("" + voltage)
    }

    return {
        sendNote: sendNote
    }
}

var api = new ArduinoApi({ outputPin: 6 })
var riff;

var scale = [0,4,7,12]//improvisr.Scales.CHROMATIC
function generateRiff() {
    riff = improvisr.generateRiff({ scale: scale })
}

var playNextNote = function() {
    var note = riff.shift()
    console.log("Playing note " + note.getChromaticOffset(scale) + " for " + note.duration)
    api.sendNote(note.getChromaticOffset(scale))

    if (! riff.length) {
        generateRiff()
    }

    setTimeout(playNextNote, note.getDurationInMs())
}

generateRiff()
playNextNote()
