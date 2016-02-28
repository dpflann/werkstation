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
        sendUdp(voltage + "," + voltage)
    }

    return {
        sendNote: sendNote
    }
}

var api = new ArduinoApi({ outputPin: 6 })
var riff, lastRiff

var config = {
    scale: improvisr.Scales.MINOR,
    numMeasures: 2,
    bpm: 180,
    loopRiff: true
}
function generateRiff() {
    riff = improvisr.generateRiff(config)
    lastRiff = _.map(riff, _.clone)
}

var playNextNote = function() {
    var note = riff.shift()
    console.log("Playing note " + note.getChromaticOffset(config.scale) + " for " + note.duration)
    api.sendNote(note.getChromaticOffset(config.scale))

    if (! riff.length) {
        if (config.loopRiff) {
            riff = _.map(lastRiff, _.clone)
        } else {
            generateRiff()
        }
    }

    setTimeout(playNextNote, note.getDurationInMs(config.bpm))
}

generateRiff()
playNextNote()
