const fs = require('fs');
const http = require('http');
const url = require('url');

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

function getQueryParams(context) {
    const urlParts = url.parse(context.request.url, true);
    return urlParts.query;
}

function getDuration(startTime) {
    const end = process.hrtime(startTime);
    return end[1] + (end[0] * 1e9);
}

function getLevel(event) {
    const intensityLevel = event.level ? parseInt(event.level) : null;
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

module.exports = async (context) => {
    const query = getQueryParams(context);

    const startTime = process.hrtime();
    const params = getParameters(query);
    if (params.error) {
        return {
            status: 200,
            body: { error: params.error },
        };
    }

    await runTest(params);

    const reused = isWarm();
    const duration = getDuration(startTime);

    return {
        status: 200,
        body: { reused, duration },
    };
};
