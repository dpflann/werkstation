var _ = require('underscore');

var Scales = {
    CHROMATIC: [0,1,2,3,4,5,6,7,8,9,10,11,12],
    MAJOR: [0,2,4,5,7,9,11,12]
}

var BPM = 120 // TODO pass dat

var Duration = {
    // TODO add rests?
    SIXTEENTH: { value: 1/16, weight: 2 },
    EIGHTH: { value: 1/8, weight: 4 },
    QUARTER: { value: 1/4, weight: 8 },
    HALF: { value: 1/2, weight: 2 },
    WHOLE: { value: 1, weight: 0 } // these are boring probably
}, weightedDurations = weightArray(Duration)

function Note(offset, duration) {
    this.offset = offset
    this.duration = duration

    this.getDurationInMs = function() {
        return 1 / BPM * 60000 * this.duration * 4
    }
    this.getChromaticOffset = function(scale) {
        return scale[this.offset]
    }
}

function generateRiff(config) {
    var measuresPerRiff = 4, // maybe config this?
        scale = config.scale

    var notes = [],
        noteOffset = 0

    while (measuresPerRiff > 0) {
        var duration = randomDuration(measuresPerRiff)
        notes.push(new Note(noteOffset, duration.value))

        measuresPerRiff -= duration.value
        noteOffset = (noteOffset + randomNoteOffsetIncrement().offset + scale.length) % scale.length

    }

    return notes
}

var OFFSETS = [
    // TODO weights should prob change according to where in the riff you are.
    { offset: -3, weight: 1 },
    { offset: -2, weight: 2 },
    { offset: -1, weight: 4 },
    { offset: 0, weight: 1 },
    { offset: 1, weight: 6 },
    { offset: 2, weight: 3 },
    { offset: 3, weight: 1 },
    { offset: 5, weight: 1 }
], weightedOffsets = weightArray(OFFSETS)

function randomNoteOffsetIncrement() {
    return _.shuffle(weightedOffsets)[0]
}

function randomDuration(max) {
    return _.chain(weightedDurations).shuffle().find(function(duration) {
        return duration.value <= max
    }).value()
}

function randomBetween(min, max) {
    return Math.floor(Math.random() * (max + 1 - min)) + min
}

function weightArray(array) {
    var weighted = []
    _.each(array, function(item) {
        for (var i = 0; i < item.weight; i++) {
            weighted.push(item)
        }
    })
    return weighted
}

module.exports = {
    Scales: Scales,
    generateRiff: generateRiff
}
