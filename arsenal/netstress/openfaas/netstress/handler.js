const fs = require('fs');
const http = require('http');

const files = {1: '/files/1Mb.dat', 2: '/files/10Mb.dat', 3: '/files/100Mb.dat'};

async function networkIntensive(level) {
    const writable = fs.createWriteStream('/dev/null');
    await new Promise((resolve) => http.get({
        host: `www.ovh.net`,
        port: 80,
        path: files[level]
    }, (res) => {
        const download = res.pipe(writable);
        download.on('close', () => resolve(res));
    }));
}

function getDuration(startTime) {
    const end = process.hrtime(startTime);
    return end[1] + (end[0] * 1e9);
}

function getLevel(event) {
    const intensityLevel = event.query.level ? parseInt(event.query.level) : null;
    if (!intensityLevel || intensityLevel < 1) {
        return {"error": "invalid level parameter"};
    }
    return intensityLevel;
}

function getParameters(event) {
    return getLevel(event);
}

async function runTest(intensityLevel){
    await networkIntensive(intensityLevel)
}

function isWarm() {
    const is_warm = process.env.warm ? true : false;
    process.env.warm = true;
    return is_warm;
}

module.exports = async (event, context) => {
    const startTime = process.hrtime();
    const params = getParameters(event);
    if (params.error) {
        return context.status(200).succeed({
            error: params.error,
        });
    }

    await runTest(params);

    const reused = isWarm();
    const duration = getDuration(startTime);

    return context.status(200).succeed({
        reused: reused,
        duration: duration,
    });
}
